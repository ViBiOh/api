package dump

import (
	"log"
	"net/http"

	"github.com/ViBiOh/httputils/pkg/dump"
)

// Handler for dump request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(dump.Request(r))
	})
}
