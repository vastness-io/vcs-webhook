package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vastness-io/queues/pkg/queue"
	"github.com/vastness-io/vcs-webhook/pkg/route/webhook"
	"github.com/vastness-io/vcs-webhook/pkg/util"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	WebhookEndpoint = "/webhook"
)

const (
	name        = "vcs-webhook"
	description = "Forwards VCS webhook event(s) to coordinator."
)

var (
	log               = logrus.WithField("pkg", "main")
	commit            string
	version           string
	addr              string
	port              int
	svcAddr           string
	debugMode         bool
	transportSecurity bool
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = description
	app.Version = fmt.Sprintf("%s (%s)", version, commit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "addr,a",
			Usage:       "TCP address to listen on",
			Value:       "127.0.0.1",
			Destination: &addr,
		},
		cli.IntFlag{
			Name:        "port,p",
			Usage:       "Port to listen on",
			Value:       8081,
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "svc,s",
			Usage:       "Service address to pass the vcs webhook payloads to.",
			Value:       "127.0.0.1:8080",
			Destination: &svcAddr,
		},
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "Debug mode",
			Destination: &debugMode,
		},
		cli.BoolTFlag{
			Name:        "transport,t",
			Usage:       "Enable Transport security",
			Destination: &transportSecurity,
		},
	}
	app.Action = func(_ *cli.Context) { run() }
	app.Run(os.Args)
}

func run() {

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	log.Info("Starting vcs-webhook")

	getClientConnection := util.NewClientConnection(svcAddr)

	var cc *grpc.ClientConn
	var err error

	if !transportSecurity {
		cc, err = getClientConnection(grpc.WithInsecure())
	} else {
		cc, err = getClientConnection()
	}

	if err != nil {
		log.Fatal(err)
	}

	q := queue.NewFIFOQueue()

	client := webhook.NewGithubWebhook(cc, q)
	defer cc.Close()

	go client.WorkOnQueue()

	r := mux.NewRouter()
	sub := r.PathPrefix(WebhookEndpoint).Subrouter()
	sub.HandleFunc(webhook.GithubWebhookEndpoint, client.OnPush).Methods("POST")
	http.Handle("/", r)

	srv := http.Server{
		Addr:    net.JoinHostPort(addr, strconv.Itoa(port)),
		Handler: r,
	}

	go func() {
		log.Infof("Listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			log.Info("Stopping workers on queue")
			q.ShutDown()
			log.Info("Exiting vcs-webhook")
			ctx, _ := context.WithTimeout(context.Background(), time.Minute)
			if err := srv.Shutdown(ctx); err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

}
