package hello

import (
	"github.com/ViBiOh/go-api/jsonHttp"
	"html"
	"net/http"
)

const delayInSeconds = 1

type hello struct {
	Name string `json:"greeting"`
}

func pluralize(s string, n int) string {
	if n > 1 {
		return (s + `s`)
	}
	return s
}

// Handler for Hello request. Should be use with net/http
type Handler struct {
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(`Access-Control-Allow-Origin`, `*`)
	w.Header().Add(`Access-Control-Allow-Headers`, `Content-Type`)
	w.Header().Add(`Access-Control-Allow-Methods`, `GET`)
	w.Header().Add(`X-Content-Type-Options`, `nosniff`)

	if r.Method == http.MethodOptions {
		w.Write(nil)
		return
	}

	hello := hello{`Hello ` + html.EscapeString(r.URL.Path) + `, I'm greeting you from the server!`}

	jsonHttp.ResponseJSON(w, hello)
}
