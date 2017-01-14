default: lint vet tst build

lint:
	go get -u github.com/golang/lint/golint
	go get -u github.com/gorilla/websocket
	golint ./...

vet:
	go vet ./...

tst:
	go test ./...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo server.go
