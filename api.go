package main

import (
	"context"
	"github.com/ViBiOh/go-api/auth"
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/healthcheck"
	"github.com/ViBiOh/go-api/hello"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"syscall"
)

const port = `1080`
const helloPath = `/hello`
const echoPath = `/echo`
const authPath = `/auth`
const healthPath = `/health`

var helloRequestMatcher = regexp.MustCompile(`^` + helloPath)
var echoRequestMatcher = regexp.MustCompile(`^` + echoPath)
var authRequestMatcher = regexp.MustCompile(`^` + authPath)
var healthRequestMatcher = regexp.MustCompile(`^` + healthPath)

func handleGracefulClose(server *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	go func() {
		<-signals
		log.Print(`SIGTERM received`)

		if server != nil {
			log.Print(`Shutting down http server`)
			if err := server.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}

		os.Exit(0)
	}()
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := []byte(r.URL.Path)

	if helloRequestMatcher.Match(urlPath) {
		http.StripPrefix(helloPath, hello.Handler{}).ServeHTTP(w, r)
	} else if echoRequestMatcher.Match(urlPath) {
		http.StripPrefix(echoPath, echo.Handler{}).ServeHTTP(w, r)
	} else if authRequestMatcher.Match(urlPath) {
		http.StripPrefix(authPath, auth.Handler{}).ServeHTTP(w, r)
	} else if healthRequestMatcher.Match(urlPath) {
		http.StripPrefix(healthPath, healthcheck.Handler{}).ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Print(`Starting server on port ` + port)

	server := &http.Server{
		Addr:    `:` + port,
		Handler: http.HandlerFunc(apiHandler),
	}

	handleGracefulClose(server)
	log.Fatal(server.ListenAndServe())
}
