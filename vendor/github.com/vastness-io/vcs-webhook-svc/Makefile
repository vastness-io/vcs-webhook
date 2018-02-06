generate:
	@echo "=> generating stubs"
	protoc -I ${PWD}/webhook/github --proto_path=${PWD}/webhook/github/ ${PWD}/webhook/github/*.proto --go_out=plugins=grpc:${PWD}/webhook/github
	@echo "=> generating mocks"
	mockgen github.com/vastness-io/vcs-webhook-svc/webhook/github GithubWebhookClient,GithubWebhookServer > mock/webhook/github/webhook_mock.go
