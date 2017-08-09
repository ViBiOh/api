package auth

import (
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils"
)

// Handler for Hello request. Should be use with net/http
type Handler struct {
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Write(nil)
		return
	}

	httputils.Unauthorized(w, fmt.Errorf(`No auth`))
}
