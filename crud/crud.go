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

func getRequestID(w http.ResponseWriter, r *http.Request) int {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, `/`))

	if err != nil {
		httputils.BadRequest(w, fmt.Errorf(`error while parsing given id %s: %v`, strings.TrimPrefix(r.URL.Path, `/`), err))
		return -1

	}
	return id
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

func updateUser(id int, name string) *user {
	foundUser, ok := users[id]

	if ok {
		foundUser.Name = name
	}

	return foundUser
}

func deleteUser(id int) {
	delete(users, id)
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Write(nil)
		return
	}

	var requestID int
	var requestUser *user

	if r.Method == http.MethodGet {
		requestID = getRequestID(w, r)

		if requestID > 0 {
			if requestUser = getUser(requestID); requestUser == nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				httputils.ResponseJSON(w, requestUser)
			}
		}
	} else if r.Method == http.MethodPost {
		if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`error while reading body: %v`, err))
		} else if err := json.Unmarshal(bodyBytes, &requestUser); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`error while unmarshalling body: %v`, err))
		} else {
			w.WriteHeader(http.StatusCreated)
			httputils.ResponseJSON(w, createUser(requestUser.Name))
		}
	} else if r.Method == http.MethodPut {
		requestID = getRequestID(w, r)

		if requestID > 0 {
			if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
				httputils.BadRequest(w, fmt.Errorf(`error while reading body: %v`, err))
			} else if err := json.Unmarshal(bodyBytes, &requestUser); err != nil {
				httputils.BadRequest(w, fmt.Errorf(`error while unmarshalling body: %v`, err))
			} else if updatedUser := updateUser(requestID, requestUser.Name); updatedUser == nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
				httputils.ResponseJSON(w, updatedUser)
			}
		}
	} else if r.Method == http.MethodDelete {
		requestID = getRequestID(w, r)

		if requestID > 0 {
			deleteUser(requestID)
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
