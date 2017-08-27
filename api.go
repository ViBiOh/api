package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
	"strings"

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
const echoPath = `/echo`
const authPath = `/auth`
const healthcheckPath = `/health`

var helloHandler = http.StripPrefix(helloPath, hello.Handler{})
var echoHandler = http.StripPrefix(echoPath, echo.Handler{})
var authHandler = http.StripPrefix(authPath, auth.Handler{})
var healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler{})

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, helloPath) {
		helloHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, echoPath) {
		echoHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, authPath) {
		authHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, healthcheckPath) {
		healthcheckHandler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	url := flag.String(`c`, ``, `URL to check`)
	tls := flag.Bool(`tls`, false, `Serve TLS content`)
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

	if *tls {
		go log.Print(cert.ListenAndServeTLS(server))
	} else {
		go log.Print(server.ListenAndServe())
	}
	httputils.ServerGracefulClose(server, nil)
}
