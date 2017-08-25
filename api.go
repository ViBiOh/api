package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	"runtime"

	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/go-api/auth"
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/healthcheck"
	"github.com/ViBiOh/go-api/hello"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cert"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
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

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if helloRequestMatcher.MatchString(r.URL.Path) {
		helloHandler.ServeHTTP(w, r)
	} else if echoRequestMatcher.MatchString(r.URL.Path) {
		echoHandler.ServeHTTP(w, r)
	} else if authRequestMatcher.MatchString(r.URL.Path) {
		authHandler.ServeHTTP(w, r)
	} else if healthcheckRequestMatcher.MatchString(r.URL.Path) {
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
		Handler: prometheus.NewPrometheusHandler(`http`, owasp.Handler{Handler: cors.Handler{Handler: http.HandlerFunc(apiHandler)}}),
	}

	go log.Panic(cert.ListenAndServeTLS(server))
	httputils.ServerGracefulClose(server, nil)
}
