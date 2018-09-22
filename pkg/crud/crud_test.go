package crud

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ViBiOh/httputils/pkg/request"
)

func Test_getRequestID(t *testing.T) {
	var cases = []struct {
		intention string
		path      string
		want      string
	}{
		{
			`should handle given string`,
			`/2`,
			`2`,
		},
	}

	for _, testCase := range cases {
		if result := getRequestID(testCase.path); result != testCase.want {
			t.Errorf("%v\ngetRequestID(%+v) = (%+v), want (%+v)", testCase.intention, testCase.path, result, testCase.want)
		}
	}
}

func Test_readCrudFromBody(t *testing.T) {
	var cases = []struct {
		intention string
		request   *http.Request
		want      *user
		wantErr   error
	}{
		{
			`should handle invalid JSON`,
			httptest.NewRequest(http.MethodGet, `/`, strings.NewReader(`{"name":"test"`)),
			nil,
			errors.New(`Error while unmarshalling body: unexpected end of JSON input`),
		},
		{
			`should handle valid JSON`,
			httptest.NewRequest(http.MethodGet, `/`, strings.NewReader(`{"id":"0","name":"test"}`)),
			&user{ID: `0`, Name: `test`},
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		result, err := readCrudFromBody(testCase.request)

		failed = false

		if err == nil && testCase.wantErr != nil {
			failed = true
		} else if err != nil && testCase.wantErr == nil {
			failed = true
		} else if err != nil && err.Error() != testCase.wantErr.Error() {
			failed = true
		} else if !reflect.DeepEqual(result, testCase.want) {
			failed = true
		}

		if failed {
			t.Errorf("%v\nreadCrudFromBody(%+v) = (%+v, %+v), want (%+v, %+v)", testCase.intention, testCase.request, result, err, testCase.want, testCase.wantErr)
		}
	}
}

func Test_listCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[string]*user
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should handle bad parsing`,
			nil,
			httptest.NewRequest(http.MethodGet, `/?page=invalid`, nil),
			`Error while parsing pagination: Error while parsing page param: strconv.ParseUint: parsing "invalid": invalid syntax
`,
			http.StatusBadRequest,
		},
		{
			`should work with empty params`,
			nil,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`{"results":[]}`,
			http.StatusOK,
		},
		{
			`should consider given args`,
			map[string]*user{`1`: {ID: `1`, Name: `1`}, `2`: {ID: `2`, Name: `2`}},
			httptest.NewRequest(http.MethodGet, `/?page=2&pageSize=1`, nil),
			`{"results":[{"id":"2","name":"2"}]}`,
			http.StatusOK,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()
		users = testCase.init

		listCrud(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nlistCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := request.ReadBodyResponse(writer.Result()); string(result) != testCase.want {
			t.Errorf("%v\nlistCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_readCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[string]*user
		request    *http.Request
		id         string
		want       string
		wantStatus int
	}{
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, nil),
			`8000`,
			`¯\_(ツ)_/¯
`,
			http.StatusNotFound,
		},
		{
			`should return serialized instance`,
			map[string]*user{`1`: {ID: `1`, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			`1`,
			`{"id":"1","name":"test"}`,
			http.StatusOK,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init
		readCrud(writer, testCase.request, testCase.id)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nreadCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := request.ReadBodyResponse(writer.Result()); string(result) != testCase.want {
			t.Errorf("%v\nreadCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_createCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should handle invalid JSON`,
			httptest.NewRequest(http.MethodGet, `/`, strings.NewReader(`{"name":"test"`)),
			`Error while parsing body: Error while unmarshalling body: unexpected end of JSON input
`,
			http.StatusBadRequest,
		},
		{
			`should create new user`,
			httptest.NewRequest(http.MethodGet, `/`, strings.NewReader(`{"name":"test"}`)),
			`"name":"test"}`,
			http.StatusCreated,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		createCrud(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\ncreateCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := request.ReadBodyResponse(writer.Result()); !strings.Contains(string(result), testCase.want) {
			t.Errorf("%v\ncreateCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_updateCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[string]*user
		request    *http.Request
		id         string
		want       string
		wantStatus int
	}{
		{
			`should handle invalid JSON`,
			nil,
			httptest.NewRequest(http.MethodGet, `/1`, strings.NewReader(`{"name":"test"`)),
			`1`,
			`Error while parsing body: Error while unmarshalling body: unexpected end of JSON input
`,
			http.StatusBadRequest,
		},
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, strings.NewReader(`{"name":"Updated Test"}`)),
			`8000`,
			`¯\_(ツ)_/¯
`,
			http.StatusNotFound,
		},
		{
			`should update given user`,
			map[string]*user{`1`: {ID: `1`, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, strings.NewReader(`{"name":"Updated Test"}`)),
			`1`,
			`{"id":"1","name":"Updated Test"}`,
			http.StatusOK,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init
		updateCrud(writer, testCase.request, testCase.id)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nupdateCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := request.ReadBodyResponse(writer.Result()); string(result) != testCase.want {
			t.Errorf("%v\nupdateCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_removeCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[string]*user
		request    *http.Request
		id         string
		want       string
		wantStatus int
	}{
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, nil),
			`8000`,
			`¯\_(ツ)_/¯
`,
			http.StatusNotFound,
		},
		{
			`should delete given user`,
			map[string]*user{`1`: {ID: `1`, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			`1`,
			``,
			http.StatusNoContent,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init
		removeCrud(writer, testCase.request, testCase.id)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nremoveCrudTest_removeCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := request.ReadBodyResponse(writer.Result()); string(result) != testCase.want {
			t.Errorf("%v\nremoveCrudTest_removeCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_ServeHTTP(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[string]*user
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should handle OPTIONS request for CORS`,
			nil,
			httptest.NewRequest(http.MethodOptions, `/`, nil),
			``,
			http.StatusNoContent,
		},
		{
			`should handle create request`,
			map[string]*user{},
			httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(`{"name":"test"}`)),
			`"name":"test"}`,
			http.StatusCreated,
		},
		{
			`should handle list request`,
			map[string]*user{`1`: {ID: `1`, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`{"results":[{"id":"1","name":"test"}]}`,
			http.StatusOK,
		},
		{
			`should handle unexpected method on root`,
			nil,
			httptest.NewRequest(http.MethodTrace, `/`, nil),
			``,
			http.StatusMethodNotAllowed,
		},
		{
			`should handle get request`,
			map[string]*user{`1`: {ID: `1`, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			`{"id":"1","name":"test"}`,
			http.StatusOK,
		},
		{
			`should handle update request`,
			map[string]*user{`1`: {ID: `1`, Name: `test`}},
			httptest.NewRequest(http.MethodPut, `/1`, strings.NewReader(`{"name":"Updated test"}`)),
			`{"id":"1","name":"Updated test"}`,
			http.StatusOK,
		},
		{
			`should handle delete request`,
			map[string]*user{`1`: {ID: `1`, Name: `test`}},
			httptest.NewRequest(http.MethodDelete, `/1`, nil),
			``,
			http.StatusNoContent,
		},
		{
			`should handle unexpected method`,
			nil,
			httptest.NewRequest(http.MethodTrace, `/1`, nil),
			``,
			http.StatusMethodNotAllowed,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init

		Handler().ServeHTTP(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nServeHTTP(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := request.ReadBodyResponse(writer.Result()); !strings.Contains(string(result), testCase.want) {
			t.Errorf("%v\nServeHTTP(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Benchmark_ServeHTTP_list(b *testing.B) {
	handler := Handler()
	users = map[string]*user{}

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, `/`, nil))
	}
}

func Benchmark_ServeHTTP_options(b *testing.B) {
	handler := Handler()
	users = map[string]*user{}

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodOptions, `/`, nil))
	}
}

func Benchmark_ServeHTTP_create(b *testing.B) {
	handler := Handler()
	users = map[string]*user{}

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(`{"name":"test"}`)))
	}
}

func Benchmark_ServeHTTP_get(b *testing.B) {
	handler := Handler()
	users = map[string]*user{`1`: {ID: `1`, Name: `test`}}

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, `/1`, nil))
	}
}
