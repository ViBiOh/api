package hello

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/httputils"
)

const locationStr = `Europe/Paris`

var location *time.Location

func init() {
	loc, err := time.LoadLocation(locationStr)
	if err != nil {
		log.Panicf(`Error while loading location %s: %v`, locationStr, err)
		return
	}

	location = loc
}

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

	httputils.ResponseJSON(w, hello{fmt.Sprintf(`Hello %s, current time is %v !`, name, time.Now().In(location))})
}
