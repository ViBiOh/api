package hello

import (
	"html"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils"
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
	if r.Method == http.MethodOptions {
		w.Write(nil)
		return
	}

	name := strings.TrimPrefix(html.EscapeString(r.URL.Path), `/`)
	if name == `` {
		name = `World`
	}

	hello := hello{`Hello ` + name + `, I'm greeting you from the server!`}

	httputils.ResponseJSON(w, hello)
}
