BIN_NAME = k3supdater

default: setup

setup:
	go install github.com/golang/mock/mockgen@latest

build:
	go build -o ./${BIN_NAME}  .

generate-mock:
	@go generate -v ./...