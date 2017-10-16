package hello

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/httputils"
)

var locationName = flag.String(`location`, `Europe/Paris`, `TimeZone for displaying current time`)
var location *time.Location

// Init hello handler
func Init() error {
	loc, err := time.LoadLocation(*locationName)
	if err != nil {
		return fmt.Errorf(`Error while loading location %s: %v`, *locationName, err)
	}

	location = loc
	return nil
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
	} else {
		name := strings.TrimPrefix(html.EscapeString(r.URL.Path), `/`)
		if name == `` {
			name = `World`
		}

		httputils.ResponseJSON(w, http.StatusOK, hello{fmt.Sprintf(`Hello %s, current time is %v !`, name, time.Now().In(location))})
	}
}
