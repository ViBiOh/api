package main

import (
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/ViBiOh/alcotest/healthcheck"
	"github.com/ViBiOh/go-api/crud"
	"github.com/ViBiOh/go-api/echo"
	"github.com/ViBiOh/go-api/hello"
	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/cors"
	"github.com/ViBiOh/httputils/owasp"
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

func main() {
	owaspConfig := owasp.Flags(``)
	corsConfig := cors.Flags(`cors`)
	helloConfig := hello.Flags(``)

	httputils.StartMainServer(func() http.Handler {
		echoHandler = http.StripPrefix(echoPath, echo.Handler())
		helloHandler = http.StripPrefix(helloPath, gziphandler.GzipHandler(hello.Handler(helloConfig)))
		crudHandler = http.StripPrefix(crudPath, gziphandler.GzipHandler(crud.Handler()))
		healthcheckHandler = http.StripPrefix(healthcheckPath, healthcheck.Handler())
		restHandler = owasp.Handler(owaspConfig, cors.Handler(corsConfig, handler()))

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, echoPath) {
				echoHandler.ServeHTTP(w, r)
			} else {
				restHandler.ServeHTTP(w, r)
			}
		})
	})
}
