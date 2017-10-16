package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/go-api/crud"
	"github.com/ViBiOh/go-api/healthcheck"
	"github.com/ViBiOh/go-api/hello"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cert"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
	"github.com/ViBiOh/httputils/rate"
)

const port = `1080`

const helloPath = `/hello`
const crudPath = `/crud`
const healthcheckPath = `/health`

var helloHandler = http.StripPrefix(helloPath, gziphandler.GzipHandler(hello.Handler{}))
var crudHandler = http.StripPrefix(crudPath, gziphandler.GzipHandler(crud.Handler{}))
var healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler{})
var restHandler = rate.Handler{Handler: owasp.Handler{Handler: cors.Handler{Handler: http.HandlerFunc(apiHandler)}}}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, healthcheckPath) {
		healthcheckHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, helloPath) {
		helloHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, crudPath) {
		crudHandler.ServeHTTP(w, r)
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

	log.Print(`Starting server on port ` + port)

	if err := hello.Init(); err != nil {
		log.Printf(`Error while initializing hello Handler: %v`, err)
	}

	server := &http.Server{
		Addr:    `:` + port,
		Handler: prometheus.NewPrometheusHandler(`http`, restHandler),
	}

	var serveError = make(chan error)
	go func() {
		defer close(serveError)
		if *tls {
			log.Print(`Listening with TLS enabled`)
			serveError <- cert.ListenAndServeTLS(server)
		} else {
			serveError <- server.ListenAndServe()
		}
	}()

	httputils.ServerGracefulClose(server, serveError, nil)
}
