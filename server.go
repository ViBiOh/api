package main

import (
	"github.com/ViBiOh/go-api/auth"
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/healthcheck"
	"github.com/ViBiOh/go-api/hello"
	"log"
	"net/http"
	"runtime"
)

const port = `1080`
const helloPath = `/hello/`
const echoPath = `/echo/`
const authPath = `/auth/`
const healthPath = `/health/`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle(helloPath, http.StripPrefix(helloPath, hello.Handler{}))
	http.Handle(echoPath, http.StripPrefix(echoPath, echo.Handler{}))
	http.Handle(authPath, http.StripPrefix(authPath, auth.Handler{}))
	http.Handle(healthPath, http.StripPrefix(healthPath, healthcheck.Handler{}))

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
