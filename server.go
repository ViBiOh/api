package main

import (
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/hello"
	"log"
	"net/http"
	"runtime"
)

const port = `1080`
const helloPath = `/hello/`
const echoPath = `/echo/`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle(helloPath, http.StripPrefix(helloPath, hello.Handler{}))
	http.Handle(echoPath, http.StripPrefix(echoPath, echo.Handler{}))

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
