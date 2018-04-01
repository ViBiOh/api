package crud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
	"github.com/ViBiOh/httputils/pkg/pagination"
	"github.com/ViBiOh/httputils/pkg/request"
)

const defaultPage = uint(1)
const defaultPageSize = uint(20)
const maxPageSize = uint(^uint(0) >> 1)

func getRequestID(path string) (uint, error) {
	parsed, err := strconv.ParseUint(strings.TrimPrefix(path, `/`), 10, 32)
	return uint(parsed), err
}

func readCrudFromBody(r *http.Request) (*user, error) {
	var requestObj user

	if bodyBytes, err := request.ReadBody(r.Body); err != nil {
		return nil, fmt.Errorf(`Error while reading body: %v`, err)
	} else if err := json.Unmarshal(bodyBytes, &requestObj); err != nil {
		return nil, fmt.Errorf(`Error while unmarshalling body: %v`, err)
	}

	return &requestObj, nil
}

func listCrud(w http.ResponseWriter, r *http.Request) {
	page, pageSize, _, _, err := pagination.ParsePaginationParams(r, defaultPageSize, maxPageSize)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf(`Error while parsing pagination: %v`, err))
		return
	}

	if err := httpjson.ResponseArrayJSON(w, http.StatusOK, listUser(page, pageSize, sortByID), httpjson.IsPretty(r.URL.RawQuery)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func readCrud(w http.ResponseWriter, r *http.Request, id uint) {
	if requestUser := getUser(id); requestUser == nil {
		httperror.NotFound(w)
	} else if err := httpjson.ResponseJSON(w, http.StatusOK, requestUser, httpjson.IsPretty(r.URL.RawQuery)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func createCrud(w http.ResponseWriter, r *http.Request) {
	if obj, err := readCrudFromBody(r); err != nil {
		httperror.BadRequest(w, fmt.Errorf(`Error while parsing body: %v`, err))
	} else if err := httpjson.ResponseJSON(w, http.StatusCreated, createUser(obj.Name), httpjson.IsPretty(r.URL.RawQuery)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func updateCrud(w http.ResponseWriter, r *http.Request, id uint) {
	if obj, err := readCrudFromBody(r); err != nil {
		httperror.BadRequest(w, fmt.Errorf(`Error while parsing body: %v`, err))
	} else if updatedUser := updateUser(id, obj.Name); updatedUser == nil {
		httperror.NotFound(w)
	} else if err := httpjson.ResponseJSON(w, http.StatusOK, updatedUser, httpjson.IsPretty(r.URL.RawQuery)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func removeCrud(w http.ResponseWriter, r *http.Request, id uint) {
	if getUser(id) == nil {
		httperror.NotFound(w)
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
				httperror.BadRequest(w, fmt.Errorf(`Error while parsing request path: %v`, err))
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
