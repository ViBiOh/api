package main

import "net/http"
import "log"
import "./hello"

const port = `1080`

func main() {
	http.HandleFunc(`/hello/`, hello.Handler)

	log.Print(`Starting server on port ` + port)
	log.Fatal(http.ListenAndServe(`:`+port, nil))
}
