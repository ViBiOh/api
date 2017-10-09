package crud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ViBiOh/httputils"
)

func Test_getRequestID(t *testing.T) {
	var cases = []struct {
		intention string
		request   *http.Request
		want      int64
		wantErr   error
	}{
		{
			`should handle empty path`,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			0,
			fmt.Errorf(`strconv.ParseInt: parsing "": invalid syntax`),
		},
		{
			`should handle invalid number`,
			httptest.NewRequest(http.MethodGet, `/abc123`, nil),
			0,
			fmt.Errorf(`strconv.ParseInt: parsing "abc123": invalid syntax`),
		},
		{
			`should handle positive number`,
			httptest.NewRequest(http.MethodGet, `/2`, nil),
			2,
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		result, err := getRequestID(testCase.request)

		failed = false

		if err == nil && testCase.wantErr != nil {
			failed = true
		} else if err != nil && testCase.wantErr == nil {
			failed = true
		} else if err != nil && err.Error() != testCase.wantErr.Error() {
			failed = true
		} else if result != testCase.want {
			failed = true
		}

		if failed {
			t.Errorf("%v\ngetRequestID(%v) = (%v, %v), want (%v, %v)", testCase.intention, testCase.request, result, err, testCase.want, testCase.wantErr)
		}
	}
}

func Test_getCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[int64]*user
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should return bad request if invalid id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`Error while parsing request id: strconv.ParseInt: parsing "": invalid syntax
`,
			http.StatusBadRequest,
		},
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, nil),
			``,
			http.StatusNotFound,
		},
		{
			`should return serialized instance`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			`{"id":1,"name":"test"}`,
			http.StatusOK,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init
		getCrud(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\ngetCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\ngetCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
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
			`Error while unmarshalling body: unexpected end of JSON input
`,
			http.StatusBadRequest,
		},
		{
			`should create new user`,
			httptest.NewRequest(http.MethodGet, `/`, strings.NewReader(`{"name":"test"}`)),
			`{"id":1,"name":"test"}`,
			http.StatusCreated,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		createCrud(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\ncreateCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\ncreateCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_updateCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[int64]*user
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should handle invalid ID`,
			nil,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`Error while parsing request id: strconv.ParseInt: parsing "": invalid syntax
`,
			http.StatusBadRequest,
		},
		{
			`should handle invalid JSON`,
			nil,
			httptest.NewRequest(http.MethodGet, `/`, strings.NewReader(`{"name":"test"`)),
			`Error while parsing request id: strconv.ParseInt: parsing "": invalid syntax
`,
			http.StatusBadRequest,
		},
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, strings.NewReader(`{"name":"Updated Test"}`)),
			``,
			http.StatusNotFound,
		},
		{
			`should update given user`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, strings.NewReader(`{"name":"Updated Test"}`)),
			`{"id":1,"name":"Updated Test"}`,
			http.StatusOK,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init
		updateCrud(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nupdateCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\nupdateCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_deleteCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[int64]*user
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should handle invalid ID`,
			nil,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`Error while parsing request id: strconv.ParseInt: parsing "": invalid syntax
`,
			http.StatusBadRequest,
		},
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, nil),
			``,
			http.StatusNotFound,
		},
		{
			`should delete given user`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			``,
			http.StatusNoContent,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init
		deleteCrud(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\ndeleteCrudTest_deleteCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\ndeleteCrudTest_deleteCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}
