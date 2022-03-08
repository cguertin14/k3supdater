BIN_NAME = k3supdater

default: setup

setup:
	@go install github.com/golang/mock/mockgen@latest

build:
	@go build -o ./${BIN_NAME}  .

test:
	@go test -v ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html

generate-mock:
	@go generate -v ./...