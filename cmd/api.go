package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/ViBiOh/go-api/pkg/dump"
	"github.com/ViBiOh/go-api/pkg/echo"
	"github.com/ViBiOh/go-api/pkg/hello"
	"github.com/ViBiOh/go-api/pkg/user"
	"github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/gzip"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/prometheus"
	"github.com/ViBiOh/httputils/pkg/rollbar"
	"github.com/ViBiOh/httputils/pkg/server"
)

const (
	echoPath  = `/echo`
	helloPath = `/hello`
	dumpPath  = `/dump`
	crudPath  = `/crud`
)

func main() {
	serverConfig := httputils.Flags(``)
	alcotestConfig := alcotest.Flags(``)
	opentracingConfig := opentracing.Flags(`tracing`)
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	prometheusConfig := prometheus.Flags(`prometheus`)
	rollbarConfig := rollbar.Flags(`rollbar`)

	helloConfig := hello.Flags(``)
	crudConfig := crud.Flags(`crud`)
	flag.Parse()

	alcotest.DoAndExit(alcotestConfig)

	serverApp := httputils.NewApp(serverConfig)
	healthcheckApp := healthcheck.NewApp()
	opentracingApp := opentracing.NewApp(opentracingConfig)
	owaspApp := owasp.NewApp(owaspConfig)
	corsApp := cors.NewApp(corsConfig)
	prometheusApp := prometheus.NewApp(prometheusConfig)
	rollbarApp := rollbar.NewApp(rollbarConfig)
	gzipApp := gzip.NewApp()

	crudApp := crud.NewApp(crudConfig, user.NewService())

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

	restHandler := server.ChainMiddlewares(handler, prometheusApp, opentracingApp, rollbarApp, gzipApp, owaspApp, corsApp)

	serverApp.ListenAndServe(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, echoPath) {
			echoHandler.ServeHTTP(w, r)
		} else {
			restHandler.ServeHTTP(w, r)
		}
	}), nil, healthcheckApp)
}
