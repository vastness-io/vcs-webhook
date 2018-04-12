generate:
	@echo "=> generating stubs"
	protoc -I ${PWD}/project --proto_path=${PWD}/project ${PWD}/project/*.proto --go_out=plugins=grpc:${PWD}/project
