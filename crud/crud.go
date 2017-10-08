package crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
)

func getRequestID(w http.ResponseWriter, r *http.Request) int64 {
	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, `/`), 10, 64)

	if err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing given id %s: %v`, strings.TrimPrefix(r.URL.Path, `/`), err))
		return -1
	}

	if id < 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	return id
}

func getCrud(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(w, r)

	if requestID > 0 {
		if requestUser := getUser(requestID); requestUser == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			httputils.ResponseJSON(w, requestUser)
		}
	}
}

func createCrud(w http.ResponseWriter, r *http.Request) {
	var requestUser *user

	if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while reading body: %v`, err))
	} else if err := json.Unmarshal(bodyBytes, &requestUser); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while unmarshalling body: %v`, err))
	} else {
		w.WriteHeader(http.StatusCreated)
		httputils.ResponseJSON(w, createUser(requestUser.Name))
	}
}

func updateCrud(w http.ResponseWriter, r *http.Request) {
	var requestUser *user
	requestID := getRequestID(w, r)

	if requestID > 0 {
		if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`Error while reading body: %v`, err))
		} else if err := json.Unmarshal(bodyBytes, &requestUser); err != nil {
			httputils.BadRequest(w, fmt.Errorf(`Error while unmarshalling body: %v`, err))
		} else if updatedUser := updateUser(requestID, requestUser.Name); updatedUser == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			httputils.ResponseJSON(w, updatedUser)
		}
	}
}

func deleteCrud(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(w, r)

	if requestID > 0 {
		deleteUser(requestID)
		w.WriteHeader(http.StatusNoContent)
	}
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
	} else if r.Method == http.MethodPost && r.URL.Path == `/` {
		createCrud(w, r)
	} else if r.Method == http.MethodGet {
		getCrud(w, r)
	} else if r.Method == http.MethodPut {
		updateCrud(w, r)
	} else if r.Method == http.MethodDelete {
		deleteCrud(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
