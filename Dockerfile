FROM alpine:3.6

RUN apk add --no-cache ca-certificates

ADD bin/linux/amd64/vcs-webhook /vcs-webhook

EXPOSE 8081

ENTRYPOINT ["/vcs-webhook"]
