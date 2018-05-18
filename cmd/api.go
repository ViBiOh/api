package main

import (
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/go-api/pkg/crud"
	"github.com/ViBiOh/go-api/pkg/dump"
	"github.com/ViBiOh/go-api/pkg/echo"
	"github.com/ViBiOh/go-api/pkg/hello"
	"github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/cors"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
)

const (
	echoPath  = `/echo`
	helloPath = `/hello`
	dumpPath  = `/dump`
	crudPath  = `/crud`
)

func main() {
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	helloConfig := hello.Flags(``)
	opentracingConfig := opentracing.Flags(`tracing`)

	healthcheckApp := healthcheck.NewApp()

	httputils.NewApp(httputils.Flags(``), func() http.Handler {
		echoHandler := http.StripPrefix(echoPath, echo.Handler())
		helloHandler := http.StripPrefix(helloPath, gziphandler.GzipHandler(hello.Handler(helloConfig)))
		dumpHandler := http.StripPrefix(dumpPath, dump.Handler())
		crudHandler := http.StripPrefix(crudPath, gziphandler.GzipHandler(crud.Handler()))

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

		healthcheckHandler := healthcheckApp.Handler(nil)
		restHandler := opentracing.NewApp(opentracingConfig).Handler(owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler)))

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == `/health` {
				healthcheckHandler.ServeHTTP(w, r)
			} else if strings.HasPrefix(r.URL.Path, echoPath) {
				echoHandler.ServeHTTP(w, r)
			} else {
				restHandler.ServeHTTP(w, r)
			}
		})
	}, nil, healthcheckApp).ListenAndServe()
}
