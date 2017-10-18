package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/alcotest"
	"github.com/ViBiOh/go-api/crud"
	"github.com/ViBiOh/go-api/echo"
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

const echoPath = `/echo`
const helloPath = `/hello`
const crudPath = `/crud`
const healthcheckPath = `/health`

var echoHandler = http.StripPrefix(echoPath, echo.Handler())
var helloHandler = http.StripPrefix(helloPath, gziphandler.GzipHandler(hello.Handler()))
var crudHandler = http.StripPrefix(crudPath, gziphandler.GzipHandler(crud.Handler()))
var healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler())
var restHandler = prometheus.Handler(`http`, rate.Handler(owasp.Handler(cors.Handler(handler()))))

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, healthcheckPath) {
			healthcheckHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, helloPath) {
			helloHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, crudPath) {
			crudHandler.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
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
		Handler: apiHandler(),
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
