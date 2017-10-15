package hello

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/httputils"
)

const locationStr = `Europe/Paris`

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

	location, err := time.LoadLocation(locationStr)
	if err != nil {
		httputils.InternalServer(w, fmt.Errorf(`Error while loading location %s: %v`, locationStr, err))
		return
	}

	httputils.ResponseJSON(w, hello{fmt.Sprintf(`Hello %s, it's %v !`, name, time.Now().In(location))})
}
