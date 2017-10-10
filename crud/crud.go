package crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
)

func getRequestID(r *http.Request) (int64, error) {
	return strconv.ParseInt(strings.TrimPrefix(r.URL.Path, `/`), 10, 64)
}

func getCrud(w http.ResponseWriter, r *http.Request) {
	if requestID, err := getRequestID(r); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing request id: %v`, err))
	} else if requestUser := getUser(requestID); requestUser == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		httputils.ResponseJSON(w, requestUser)
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

	if requestID, err := getRequestID(r); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing request id: %v`, err))
	} else if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
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

func deleteCrud(w http.ResponseWriter, r *http.Request) {
	if requestID, err := getRequestID(r); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing request id: %v`, err))
	} else if getUser(requestID) == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		deleteUser(requestID)
		w.WriteHeader(http.StatusNoContent)
	}
}

// Handler for CRUD request. Should be use with net/http
type Handler struct {
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
