package crud

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = map[int]*user{}

// Handler for Hello request. Should be use with net/http
type Handler struct {
}

func get(id int) *user {
	return users[id]
}

func create(created user) int {
	created.ID = rand.Int()
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
			id := create(bodyUser)
			w.Write([]byte(strconv.Itoa(id)))
		}
	} else if r.Method == http.MethodGet {
		if id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, `/`)); err != nil {
			httputils.BadRequest(w, err)
		} else if foundUser := get(id); foundUser != nil {
			httputils.ResponseJSON(w, foundUser)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
