package crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils"
	"github.com/ViBiOh/httputils/pagination"
)

const defaultPage = uint(1)
const defaultPageSize = uint(20)
const maxPageSize = uint(^uint(0) >> 1)

func getRequestID(path string) (uint, error) {
	parsed, err := strconv.ParseUint(strings.TrimPrefix(path, `/`), 10, 32)
	return uint(parsed), err
}

func listCrud(w http.ResponseWriter, r *http.Request) {
	page, pageSize, _, _, err := pagination.ParsePaginationParams(r, defaultPageSize, maxPageSize)
	if err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while parsing pagination: %v`, err))
		return
	}

	httputils.ResponseArrayJSON(w, http.StatusOK, listUser(page, pageSize, sortByID), httputils.IsPretty(r.URL.RawQuery))
}

func readCrud(w http.ResponseWriter, r *http.Request, id uint) {
	if requestUser := getUser(id); requestUser == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		httputils.ResponseJSON(w, http.StatusOK, requestUser, httputils.IsPretty(r.URL.RawQuery))
	}
}

func createCrud(w http.ResponseWriter, r *http.Request) {
	var requestUser *user

	if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while reading body: %v`, err))
	} else if err := json.Unmarshal(bodyBytes, &requestUser); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while unmarshalling body: %v`, err))
	} else {
		httputils.ResponseJSON(w, http.StatusCreated, createUser(requestUser.Name), httputils.IsPretty(r.URL.RawQuery))
	}
}

func updateCrud(w http.ResponseWriter, r *http.Request, id uint) {
	var requestUser *user

	if bodyBytes, err := httputils.ReadBody(r.Body); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while reading body: %v`, err))
	} else if err := json.Unmarshal(bodyBytes, &requestUser); err != nil {
		httputils.BadRequest(w, fmt.Errorf(`Error while unmarshalling body: %v`, err))
	} else if updatedUser := updateUser(id, requestUser.Name); updatedUser == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		httputils.ResponseJSON(w, http.StatusOK, updatedUser, httputils.IsPretty(r.URL.RawQuery))
	}
}

func removeCrud(w http.ResponseWriter, r *http.Request, id uint) {
	if getUser(id) == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		deleteUser(id)
		w.WriteHeader(http.StatusNoContent)
	}
}

// Handler for CRUD request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
		} else if r.URL.Path == `/` || r.URL.Path == `` {
			if r.Method == http.MethodPost {
				createCrud(w, r)
			} else if r.Method == http.MethodGet {
				listCrud(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {
			if id, err := getRequestID(r.URL.Path); err != nil {
				httputils.BadRequest(w, fmt.Errorf(`Error while parsing request path: %v`, err))
			} else if r.Method == http.MethodGet {
				readCrud(w, r, id)
			} else if r.Method == http.MethodPut {
				updateCrud(w, r, id)
			} else if r.Method == http.MethodDelete {
				removeCrud(w, r, id)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	})
}
