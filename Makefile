default: go docker

go: deps dev

dev: format lint tst bench build

docker: docker-deps docker-build

deps:
	go get -t ./...
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/goimports

format:
	goimports -w **/*.go *.go
	gofmt -s -w **/*.go *.go

lint:
	golint ./...
	go vet ./...

tst:
	script/coverage

bench:
	go test ./... -bench . -benchmem -run Benchmark.*

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/api api.go

docker-deps:
	curl -s -o cacert.pem https://curl.haxx.se/ca/cacert.pem
	curl -s -o zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

docker-build:
	docker build -t ${DOCKER_USER}/api .

docker-push:
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}
	docker push ${DOCKER_USER}/api
