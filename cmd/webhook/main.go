package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vastness-io/queues/pkg/queue"
	"google.golang.org/grpc"

	toolkit "github.com/vastness-io/toolkit/pkg/grpc"
	"github.com/vastness-io/vcs-webhook-svc/webhook"
	"github.com/vastness-io/vcs-webhook/pkg/route"
	"github.com/vastness-io/vcs-webhook/pkg/service"
	"github.com/vastness-io/vcs-webhook/pkg/transport"
)

const (
	name        = "vcs-webhook"
	description = "Forwards VCS webhook event(s) to coordinator."
)

var (
	log       = logrus.WithField("component", name)
	commit    string
	version   string
	addr      string
	port      int
	svcAddr   string
	capacity  int64
	fillRate  time.Duration
	debugMode bool
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
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
			Usage:       "Service address to pass the vcs webhook payloads to",
			Value:       "127.0.0.1:8080",
			Destination: &svcAddr,
		},
		cli.Int64Flag{
			Name:        "capacity,c",
			Usage:       "Queue capacity",
			Value:       1000,
			Destination: &capacity,
		},
		cli.DurationFlag{
			Name:        "rate,r",
			Usage:       "Token bucket fill rate, to rate limit queue processing",
			Value:       time.Millisecond * 250,
			Destination: &fillRate,
		},
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "Debug mode",
			Destination: &debugMode,
		},
	}
	app.Action = func(_ *cli.Context) { run() }
	app.Run(os.Args)
}

func run() {

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	log.Infof("Starting %s", name)

	var (
		tracer  = opentracing.GlobalTracer()
		cc, err = toolkit.NewGRPCClient(tracer, log, grpc.WithInsecure())(svcAddr)
	)

	if err != nil {
		log.Fatal(err)
	}

	var (
		vcsEventClient         = vcs.NewVcsEventClient(cc)
		githubReqQ             = queue.NewFIFOQueue(capacity, fillRate)
		bitbucketServerReqQ    = queue.NewFIFOQueue(capacity, fillRate)
		githubService          = service.NewGithubWebhookService(vcsEventClient, githubReqQ)
		bitbucketServerService = service.NewBitbucketServerWebhookService(vcsEventClient, bitbucketServerReqQ)
		githubRoute            = &route.VCSRoute{
			Service: githubService,
		}
		bitbucketServerRoute = &route.VCSRoute{
			Service: bitbucketServerService,
		}
		httpTransport = transport.NewHTTPTransport(githubRoute, bitbucketServerRoute)
		address       = net.JoinHostPort(addr, strconv.Itoa(port))
		srv           = http.Server{
			Addr:    address,
			Handler: httpTransport,
		}
	)

	defer cc.Close()

	go githubService.Work()

	go bitbucketServerService.Work()

	go func() {
		log.Infof("Listening on %s", address)
		srv.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			log.Info("Stopping workers on queue")
			githubReqQ.ShutDown()
			bitbucketServerReqQ.ShutDown()
			log.Infof("Exiting %s", name)
			ctx, _ := context.WithTimeout(context.Background(), time.Minute)
			if err := srv.Shutdown(ctx); err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

}
