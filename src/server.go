package main

import (
	"./hello"
	"log"
	"net/http"
)

const port = `1080`
const helloPath = `/hello/`

func main() {
	http.Handle(helloPath, http.StripPrefix(helloPath, hello.Handler{}))

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
