package crud

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/httpjson"
	"github.com/ViBiOh/httputils/pkg/pagination"
	"github.com/ViBiOh/httputils/pkg/request"
	"github.com/ViBiOh/httputils/pkg/tools"
)

var (
	// ErrNotFound occurs when item with given ID if not found
	ErrNotFound = errors.New(`Item not found`)
)

// App stores informations
type App struct {
	defaultPage     uint
	defaultPageSize uint
	maxPageSize     uint
	service         ItemService
}

// NewApp creates new App from Flags' config
func NewApp(config map[string]*uint, service ItemService) *App {
	return &App{
		defaultPage:     *config[`defaultPage`],
		defaultPageSize: *config[`defaultPageSize`],
		maxPageSize:     *config[`maxPageSize`],
		service:         service,
	}
}

// Flags adds flags for given prefix
func Flags(prefix string) map[string]*uint {
	return map[string]*uint{
		`defaultPage`:     flag.Uint(tools.ToCamel(fmt.Sprintf(`%sDefaultPage`, prefix)), 1, fmt.Sprintf(`[%s] Default page`, prefix)),
		`defaultPageSize`: flag.Uint(tools.ToCamel(fmt.Sprintf(`%sDefaultPageSize`, prefix)), 20, fmt.Sprintf(`[%s] Default page size`, prefix)),
		`maxPageSize`:     flag.Uint(tools.ToCamel(fmt.Sprintf(`%sMaxPageSize`, prefix)), 500, fmt.Sprintf(`[%s] Max page size`, prefix)),
	}
}

func getID(r *http.Request) string {
	return strings.Split(strings.Trim(r.URL.Path, `/`), `/`)[0]
}

func (a App) readPayload(r *http.Request) (Item, error) {
	bodyBytes, err := request.ReadBodyRequest(r)

	if err != nil {
		return nil, fmt.Errorf(`Error while reading body: %v`, err)
	}

	var obj = a.service.Empty()
	if err := json.Unmarshal(bodyBytes, obj); err != nil {
		return nil, fmt.Errorf(`Error while unmarshalling body: %v`, err)
	}

	return obj, nil
}

func (a App) list(w http.ResponseWriter, r *http.Request) {
	page, pageSize, _, _, err := pagination.ParseParams(r, a.defaultPage, a.defaultPageSize, a.maxPageSize)
	if err != nil {
		httperror.BadRequest(w, fmt.Errorf(`Error while parsing pagination: %v`, err))
		return
	}

	if err := httpjson.ResponseArrayJSON(w, http.StatusOK, a.service.List(page, pageSize), httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
	}
}

func (a App) get(w http.ResponseWriter, r *http.Request, id string) {
	obj := a.service.Get(id)

	if obj == nil {
		httperror.NotFound(w)
		return
	}

	if err := httpjson.ResponseJSON(w, http.StatusOK, obj, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (a App) create(w http.ResponseWriter, r *http.Request) {
	obj, err := a.readPayload(r)

	if err != nil {
		httperror.BadRequest(w, fmt.Errorf(`Error while parsing body: %v`, err))
		return
	}

	obj, err = a.service.Create(obj)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := httpjson.ResponseJSON(w, http.StatusCreated, obj, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (a App) update(w http.ResponseWriter, r *http.Request, id string) {
	obj, err := a.readPayload(r)

	if err != nil {
		httperror.BadRequest(w, fmt.Errorf(`Error while parsing body: %v`, err))
		return
	}

	obj, err = a.service.Update(id, obj)
	if err == ErrNotFound {
		httperror.NotFound(w)
		return
	}

	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := httpjson.ResponseJSON(w, http.StatusOK, obj, httpjson.IsPretty(r)); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (a App) delete(w http.ResponseWriter, r *http.Request, id string) {
	if a.service.Get(id) == nil {
		httperror.NotFound(w)
		return
	}

	err := a.service.Delete(id)
	if err != nil {
		if err == ErrNotFound {
			httperror.NotFound(w)
		} else {
			httperror.InternalServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handler for CRUD requests. Should be use with net/http
func (a App) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if tools.IsRoot(r) {
				a.list(w, r)
			} else {
				a.get(w, r, getID(r))
			}

		case http.MethodPost:
			if tools.IsRoot(r) {
				a.create(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}

		case http.MethodPut:
			if !tools.IsRoot(r) {
				a.update(w, r, getID(r))
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}

		case http.MethodDelete:
			if !tools.IsRoot(r) {
				a.delete(w, r, getID(r))
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}

		case http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
