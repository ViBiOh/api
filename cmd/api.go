package main

import (
	"flag"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ViBiOh/api/pkg/dump"
	"github.com/ViBiOh/api/pkg/echo"
	"github.com/ViBiOh/api/pkg/hello"
	"github.com/ViBiOh/api/pkg/user"
	httputils "github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/crud"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/prometheus"
)

const (
	echoPath  = "/echo"
	helloPath = "/hello"
	dumpPath  = "/dump"
	crudPath  = "/crud"
	docPath   = "doc/"
)

func main() {
	fs := flag.NewFlagSet("api", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	opentracingConfig := opentracing.Flags(fs, "tracing")
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

	helloConfig := hello.Flags(fs, "")
	crudConfig := crud.Flags(fs, "crud")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

	prometheusApp := prometheus.New(prometheusConfig)
	opentracingApp := opentracing.New(opentracingConfig)
	owaspApp := owasp.New(owaspConfig)
	corsApp := cors.New(corsConfig)

	crudApp := crud.New(crudConfig, user.New())

	helloHandler := http.StripPrefix(helloPath, hello.Handler(helloConfig))
	dumpHandler := http.StripPrefix(dumpPath, dump.Handler())
	echoHandler := http.StripPrefix(echoPath, echo.Handler())
	crudHandler := http.StripPrefix(crudPath, crudApp.Handler())

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, helloPath) {
			helloHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, dumpPath) {
			dumpHandler.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, crudPath) {
			crudHandler.ServeHTTP(w, r)
		} else {
			w.Header().Set("Cache-Control", "no-cache")
			http.ServeFile(w, r, path.Join(docPath, r.URL.Path))
		}
	})

	restHandler := httputils.ChainMiddlewares(handler, prometheusApp, opentracingApp, owaspApp, corsApp)

	httputils.New(serverConfig).ListenAndServe(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, echoPath) {
			echoHandler.ServeHTTP(w, r)
		} else {
			restHandler.ServeHTTP(w, r)
		}
	}), httputils.HealthHandler(nil), nil)
}
