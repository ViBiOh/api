package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/alcotest/healthcheck"
	"github.com/ViBiOh/go-api/crud"
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/hello"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cert"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/owasp"
	"github.com/ViBiOh/httputils/prometheus"
	"github.com/ViBiOh/httputils/rate"
)

const echoPath = `/echo`
const helloPath = `/hello`
const crudPath = `/crud`
const healthcheckPath = `/health`

var (
	echoHandler        http.Handler
	helloHandler       http.Handler
	crudHandler        http.Handler
	healthcheckHandler http.Handler
	restHandler        http.Handler
)

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, healthcheckPath) {
			healthcheckHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, helloPath) {
			helloHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, crudPath) {
			crudHandler.ServeHTTP(w, r)
		} else {
			httputils.NotFound(w)
		}
	})
}

func apiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, echoPath) {
			echoHandler.ServeHTTP(w, r)
		} else {
			restHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	port := flag.String(`port`, `1080`, `Listen port`)
	tls := flag.Bool(`tls`, false, `Serve TLS content`)
	alcotestConfig := alcotest.Flags(``)
	certConfig := cert.Flags(`tls`)
	prometheusConfig := prometheus.Flags(`prometheus`)
	rateConfig := rate.Flags(`rate`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	helloConfig := hello.Flags(``)
	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	log.Printf(`Starting server on port %s`, *port)

	echoHandler = http.StripPrefix(echoPath, echo.Handler())
	helloHandler = http.StripPrefix(helloPath, gziphandler.GzipHandler(hello.Handler(helloConfig)))
	crudHandler = http.StripPrefix(crudPath, gziphandler.GzipHandler(crud.Handler()))
	healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler())
	restHandler = prometheus.Handler(prometheusConfig, rate.Handler(rateConfig, owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler()))))
	server := &http.Server{
		Addr:    `:` + *port,
		Handler: apiHandler(),
	}

	var serveError = make(chan error)
	go func() {
		defer close(serveError)
		if *tls {
			log.Print(`Listening with TLS enabled`)
			serveError <- cert.ListenAndServeTLS(certConfig, server)
		} else {
			serveError <- server.ListenAndServe()
		}
	}()

	httputils.ServerGracefulClose(server, serveError, nil)
}
