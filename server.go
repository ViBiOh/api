package main

import (
	"github.com/ViBiOh/go-api/hello"
	"log"
	"net/http"
	"runtime"
)

const port = `1080`
const helloPath = `/hello/`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.Handle(helloPath, http.StripPrefix(helloPath, hello.Handler{}))

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
