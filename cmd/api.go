package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/ViBiOh/go-api/pkg/dump"
	"github.com/ViBiOh/go-api/pkg/echo"
	"github.com/ViBiOh/go-api/pkg/hello"
	"github.com/ViBiOh/go-api/pkg/user"
	httputils "github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/gzip"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/prometheus"
	"github.com/ViBiOh/httputils/pkg/server"
)

const (
	echoPath  = `/echo`
	helloPath = `/hello`
	dumpPath  = `/dump`
	crudPath  = `/crud`
)

func main() {
	fs := flag.NewFlagSet(`api`, flag.ExitOnError)

	serverConfig := httputils.Flags(fs, ``)
	alcotestConfig := alcotest.Flags(fs, ``)
	prometheusConfig := prometheus.Flags(fs, `prometheus`)
	opentracingConfig := opentracing.Flags(fs, `tracing`)
	owaspConfig := owasp.Flags(fs, ``)
	corsConfig := cors.Flags(fs, `cors`)

	helloConfig := hello.Flags(fs, ``)
	crudConfig := crud.Flags(fs, `crud`)

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Fatal(`%+v`, err)
	}

	alcotest.DoAndExit(alcotestConfig)

	serverApp := httputils.New(serverConfig)
	healthcheckApp := healthcheck.New()
	prometheusApp := prometheus.New(prometheusConfig)
	opentracingApp := opentracing.New(opentracingConfig)
	gzipApp := gzip.New()
	owaspApp := owasp.New(owaspConfig)
	corsApp := cors.New(corsConfig)

	crudApp := crud.New(crudConfig, user.New())

	helloHandler := http.StripPrefix(helloPath, hello.Handler(helloConfig))
	crudHandler := http.StripPrefix(crudPath, crudApp.Handler())
	dumpHandler := http.StripPrefix(dumpPath, dump.Handler())
	echoHandler := http.StripPrefix(echoPath, echo.Handler())

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, helloPath) {
			helloHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, dumpPath) {
			dumpHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, crudPath) {
			crudHandler.ServeHTTP(w, r)
		} else {
			httperror.NotFound(w)
		}
	})

	restHandler := server.ChainMiddlewares(handler, prometheusApp, opentracingApp, gzipApp, owaspApp, corsApp)

	serverApp.ListenAndServe(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, echoPath) {
			echoHandler.ServeHTTP(w, r)
		} else {
			restHandler.ServeHTTP(w, r)
		}
	}), nil, healthcheckApp)
}
