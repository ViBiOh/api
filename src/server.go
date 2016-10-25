package main

import "net/http"
import "html"
import "encoding/json"
import "time"
import "log"
import "strings"
import "strconv"

const delayInSeconds = 1
const port = "1080"

type Hello struct {
  Name string `json:"greeting"`
}

func pluralize(s string, n int) string {
  if n > 1 {
    return (s + "s")
  }
  return s
}

func apiHello(w http.ResponseWriter, r *http.Request) {
  time.Sleep(delayInSeconds * time.Second)
  hello := Hello{"Hello " + html.EscapeString(strings.Replace(r.URL.Path, "/api/hello/", "", -1)) + ", I'm greeting you from the server with " + strconv.Itoa(delayInSeconds) + " " + pluralize("second", delayInSeconds) + " delay"}

  helloJson, err := json.Marshal(hello)
  if err == nil {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write(helloJson)
  } else {
    log.Fatal(err)
  }
}

func main() {
  http.HandleFunc("/hello/", apiHello)

  log.Print("Starting server on port " + port)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}
