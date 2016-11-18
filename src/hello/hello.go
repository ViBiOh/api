package hello

import "net/http"
import "html"
import "strings"
import "../jsonHttp"

const delayInSeconds = 1

type Hello struct {
	Name string `json:"greeting"`
}

func pluralize(s string, n int) string {
	if n > 1 {
		return (s + `s`)
	}
	return s
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(`Access-Control-Allow-Origin`, `*`)
	w.Header().Add(`Access-Control-Allow-Headers`, `Content-Type`)
	w.Header().Add(`Access-Control-Allow-Methods`, `GET`)
	w.Header().Add(`X-Content-Type-Options`, `nosniff`)

	hello := Hello{`Hello ` + html.EscapeString(r.URL.Path) + `, I'm greeting you from the server!`}

	jsonHttp.ResponseJson(w, hello)
}
