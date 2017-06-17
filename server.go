package main

import (
	"github.com/ViBiOh/go-api/auth"
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/hello"
	"github.com/ViBiOh/go-api/status"
	"log"
	"net/http"
	"runtime"
)

const port = `1080`
const helloPath = `/hello/`
const echoPath = `/echo/`
const authPath = `/auth/`
const statusPath = `/status/`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle(helloPath, http.StripPrefix(helloPath, hello.Handler{}))
	http.Handle(echoPath, http.StripPrefix(echoPath, echo.Handler{}))
	http.Handle(authPath, http.StripPrefix(authPath, auth.Handler{}))
	http.Handle(statusPath, http.StripPrefix(statusPath, status.Handler{}))

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
