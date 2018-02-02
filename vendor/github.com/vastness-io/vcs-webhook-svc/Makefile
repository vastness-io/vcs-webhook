generate:
	@echo "=> generating stubs"
	protoc -I ${PWD}/webhook/github --proto_path=${PWD}/webhook/github/ ${PWD}/webhook/github/*.proto --go_out=plugins=grpc:${PWD}/webhook/github
