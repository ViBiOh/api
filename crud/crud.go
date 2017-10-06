package crud

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils"
	"github.com/satori/go.uuid"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = map[string]*user{}

// Handler for Hello request. Should be use with net/http
type Handler struct {
}

func get(id string) *user {
	return users[id]
}

func create(created user) string {
	created.ID = uuid.NewV4().String()
	users[created.ID] = &created

	return created.ID
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Write(nil)
		return
	}

	var bodyUser user

	if r.Method == http.MethodPost {
		if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
			httputils.BadRequest(w, err)
		} else if err := json.Unmarshal(bodyBytes, &bodyUser); err != nil {
			httputils.BadRequest(w, err)
		} else {
			w.Write([]byte(create(bodyUser)))
		}
	} else if r.Method == http.MethodGet {
		if foundUser := get(strings.TrimPrefix(r.URL.Path, `/`)); foundUser != nil {
			httputils.ResponseJSON(w, foundUser)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
