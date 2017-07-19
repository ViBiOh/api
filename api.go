package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"syscall"
	"time"

	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/go-api/auth"
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/healthcheck"
	"github.com/ViBiOh/go-api/hello"
)

const port = `1080`

const helloPath = `/hello`

var helloRequestMatcher = regexp.MustCompile(`^` + helloPath)
var helloHandler = http.StripPrefix(helloPath, hello.Handler{})

const echoPath = `/echo`

var echoRequestMatcher = regexp.MustCompile(`^` + echoPath)
var echoHandler = http.StripPrefix(echoPath, echo.Handler{})

const authPath = `/auth`

var authRequestMatcher = regexp.MustCompile(`^` + authPath)
var authHandler = http.StripPrefix(authPath, auth.Handler{})

const healthcheckPath = `/health`

var healthcheckRequestMatcher = regexp.MustCompile(`^` + healthcheckPath)
var healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler{})

func handleGracefulClose(server *http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	<-signals

	log.Print(`SIGTERM received`)

	if server != nil {
		log.Print(`Shutting down http server`)

		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		
		if err := server.Shutdown(ctx); err != nil {
			log.Print(err)
		}
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := []byte(r.URL.Path)

	if helloRequestMatcher.Match(urlPath) {
		helloHandler.ServeHTTP(w, r)
	} else if echoRequestMatcher.Match(urlPath) {
		echoHandler.ServeHTTP(w, r)
	} else if authRequestMatcher.Match(urlPath) {
		authHandler.ServeHTTP(w, r)
	} else if healthcheckRequestMatcher.Match(urlPath) {
		healthcheckHandler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	url := flag.String(`c`, ``, `URL to check`)
	flag.Parse()

	if *url != `` {
		alcotest.Do(url)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Print(`Starting server on port ` + port)

	server := &http.Server{
		Addr:    `:` + port,
		Handler: http.HandlerFunc(apiHandler),
	}

	go server.ListenAndServe()
	handleGracefulClose(server)
}
