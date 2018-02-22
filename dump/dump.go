package dump

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils/request"
)

// Handler for dump request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headers string
		for key, value := range r.Header {
			headers += fmt.Sprintf("%s: %s\n", key, strings.Join(value, `,`))
		}

		body, err := request.ReadBody(r.Body)
		if err != nil {
			log.Printf(`Error while reading body: %v`, err)
		}

		log.Printf("\n%s %s\n%s\n%s", r.Method, r.URL.Path, headers, body)
	})
}
