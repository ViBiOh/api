package main

import (
	"./hello"
	"log"
	"net/http"
)

const port = `1080`
const HELLO_PATH = `/hello/`

func main() {
	http.Handle(HELLO_PATH, http.StripPrefix(HELLO_PATH, hello.Handler{}))

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
