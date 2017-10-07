package crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users = map[int]*user{}
	seq   = 1
)

// Handler for Hello request. Should be use with net/http
type Handler struct {
}

func getUser(id int) *user {
	return users[id]
}

func createUser(name string) *user {
	createdUser := &user{ID: seq, Name: name}
	users[seq] = createdUser

	seq++
	return createdUser
}

func deleteUser(id int) {
	delete(users, id)
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Write(nil)
		return
	}

	var bodyUser user

	if r.Method == http.MethodGet {
		if id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, `/`)); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`error while parsing given id %s: %v`, strings.TrimPrefix(r.URL.Path, `/`), err))
		} else if foundUser := getUser(id); foundUser == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			httputils.ResponseJSON(w, foundUser)
		}
	} else if r.Method == http.MethodPost {
		if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`error while reading body: %v`, err))
		} else if err := json.Unmarshal(bodyBytes, &bodyUser); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`error while unmarshalling body: %v`, err))
		} else {
			w.WriteHeader(http.StatusCreated)
			httputils.ResponseJSON(w, createUser(bodyUser.Name))
		}
	} else if r.Method == http.MethodDelete {
		if id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, `/`)); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`error while parsing given id %s: %v`, strings.TrimPrefix(r.URL.Path, `/`), err))
		} else {
			deleteUser(id)
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
