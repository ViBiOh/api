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
	"github.com/ViBiOh/httputils/v3/pkg/alcotest"
	"github.com/ViBiOh/httputils/v3/pkg/cors"
	"github.com/ViBiOh/httputils/v3/pkg/crud"
	"github.com/ViBiOh/httputils/v3/pkg/httputils"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
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
	owaspConfig := owasp.Flags(fs, "")
	corsConfig := cors.Flags(fs, "cors")

	helloConfig := hello.Flags(fs, "")
	crudConfig := crud.Flags(fs, "crud")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

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

	restHandler := httputils.ChainMiddlewares(handler, prometheus.New(prometheusConfig), owasp.New(owaspConfig), cors.New(corsConfig))

	server := httputils.New(serverConfig)
	server.ListenServeWait(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, echoPath) {
			echoHandler.ServeHTTP(w, r)
		} else {
			restHandler.ServeHTTP(w, r)
		}
	}))
}
