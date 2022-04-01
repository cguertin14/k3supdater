BIN_NAME = k3supdater

default: setup

setup:
	@go install github.com/golang/mock/mockgen@latest

build:
	@go build -o ./${BIN_NAME}  .

test:
	@go test --tags unit -v ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html

generate-mock:
	@go generate -v ./...