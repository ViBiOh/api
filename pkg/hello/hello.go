package hello

import (
	"flag"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/tools"
)

// Hello represents the outputed welcome message
// swagger:model
type Hello struct {
	// required: true
	Name string `json:"greeting"`
}

// Config of package
type Config struct {
	locationName *string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		locationName: fs.String(tools.ToCamel(fmt.Sprintf("%sLocation", prefix)), "Europe/Paris", "[hello] TimeZone for displaying current time"),
	}
}

// Handler for Hello request. Should be use with net/http
func Handler(config Config) http.Handler {
	location, err := time.LoadLocation(*config.locationName)
	if err != nil {
		logger.Error("Error while loading location %s: %v", *config.locationName, err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			if _, err := w.Write(nil); err != nil {
				httperror.InternalServerError(w, err)
			}
			return
		}

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// swagger:route GET /hello hello
		//
		// Say hello to caller
		//
		// Consumes:
		// - application/json
		//
		// Produces:
		// - application/json
		//
		// Responses:
		// 200: Hello
		name := strings.TrimPrefix(html.EscapeString(r.URL.Path), "/")
		if name == "" {
			name = "World"
		}

		if err := httpjson.ResponseJSON(w, http.StatusOK, Hello{fmt.Sprintf("Hello %s, current time is %v !", name, time.Now().In(location))}, httpjson.IsPretty(r)); err != nil {
			httperror.InternalServerError(w, err)
		}
	})
}
