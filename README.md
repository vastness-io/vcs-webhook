VCS Webhook consumes VCS specific *Event* payloads and send them to the [Vastness](https://github.com/QualityAlliance/vastness) API server.  
[![Build Status](https://travis-ci.org/vastness-io/vcs-webhook.svg)](https://travis-ci.org/vastness-io/vcs-webhook)
[![GoDoc](https://godoc.org/github.com/vastness-io/vcs-webhook?status.svg)](https://godoc.org/github.com/vastness-io/vcs-webhook)
[![codecov](https://codecov.io/gh/vastness-io/vcs-webhook/branch/master/graph/badge.svg)](https://codecov.io/gh/vastness-io/vcs-webhook)
---
#### Supported VCS
* Github
    * Push Event
* Bitbucket Server
    * Post Webhook plugin


### Usage
```bash
./vcs-webhook # Will start the application with defaults
# For help use the -h flag
```

### Development
#### Prerequistities
* Go 1.9.x
* [glide](https://github.com/Masterminds/glide)
* Make

```bash
go get -d github.com/vastness-io/vcs-webhook
cd $GOPATH/src/github.com/vastness-io/vcs-webhook
make build
```
