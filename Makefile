default: deps lint tst build

deps:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/gorilla/websocket
	go get -u github.com/ViBiOh/alcotest/alcotest
	go get -u github.com/ViBiOh/httputils
	go get -u github.com/ViBiOh/httputils/prometheus

fmt:
	goimports -w **/*.go *.go
	gofmt -s -w **/*.go *.go

lint:
	golint ./...
	go vet ./...

tst:
	script/coverage

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/api api.go
