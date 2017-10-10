package hello

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils"
)

type hello struct {
	Name string `json:"greeting"`
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

	httputils.ResponseJSON(w, hello{fmt.Sprintf(`Hello %s, I'm greeting you from the server!`, name)})
}
