MAKEFLAGS += --silent
GOBIN=bin
PID=/tmp/.$(APP_NAME).pid

APP_NAME ?= api
VERSION ?= $(shell git log --pretty=format:'%h' -n 1)
AUTHOR ?= $(shell git log --pretty=format:'%an' -n 1)

help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## $(APP_NAME): Build App with dependencies download
$(APP_NAME): deps go

## go: Build App
go: format lint tst bench build

## name: Output app name
name:
	@echo -n $(APP_NAME)

## version: Output last commit sha1
version:
	@echo -n $(VERSION)

## author: Output last commit author
author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(AUTHOR)"

## deps: Download dependencies
deps:
	go get github.com/golang/dep/cmd/dep
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck
	go get golang.org/x/tools/cmd/goimports
	dep ensure

## format: Format code
format:
	goimports -w */*.go */*/*.go
	gofmt -s -w */*.go */*/*.go

## Lint: Lint code
lint:
	golint `go list ./... | grep -v vendor`
	errcheck -ignoretests `go list ./... | grep -v vendor`
	go vet ./...

## tst: Test code with coverage
tst:
	script/coverage

## bench: Benchmark code
bench:
	go test ./... -bench . -benchmem -run Benchmark.*

## build: Build binary of App
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(GOBIN)/$(APP_NAME) cmd/api.go

## start: Start app
start: stop build
	$(GOBIN)/$(APP_NAME) \
		-tls=false \
	& echo $$! > $(PID)

## stop: Stop app
stop:
	touch $(PID)
	kill -9 `cat $(PID)` 2> /dev/null || true
	rm $(PID)

.PHONY: $(APP_NAME) go name version author deps format lint tst bench build start stop
