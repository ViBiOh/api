package dump

import (
	"net/http"

	"github.com/ViBiOh/httputils/pkg/dump"
	"github.com/ViBiOh/httputils/pkg/httperror"
)

// Handler for dump request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(dump.Request(r))); err != nil {
			httperror.InternalServerError(w, err)
		}
	})
}
