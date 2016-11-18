package main

import "net/http"
import "log"
import "./hello"

const port = `1080`
const HELLO_PATH = `/hello/`

func main() {
	http.HandleFunc(HELLO_PATH, http.StripPrefix(HELLO_PATH, hello.Handler)

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
