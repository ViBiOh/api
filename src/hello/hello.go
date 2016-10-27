package hello

import "net/http"
import "html"
import "strings"
import "../json"

const delayInSeconds = 1

type Hello struct {
	Name string `json:"greeting"`
}

func pluralize(s string, n int) string {
	if n > 1 {
		return (s + "s")
	}
	return s
}

func Handler(w http.ResponseWriter, r *http.Request) {
	hello := Hello{"Hello " + html.EscapeString(strings.Replace(r.URL.Path, "/hello/", "", -1)) + ", I'm greeting you from the server!"}

	json.ResponseJson(w, hello)
}
