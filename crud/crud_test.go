package crud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ViBiOh/httputils"
)

func Test_listCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[uint]*user
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should work with empty params`,
			nil,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`{"results":[]}`,
			http.StatusOK,
		},
		{
			`should consider given args`,
			map[uint]*user{1: {ID: 1, Name: `1`}, 2: {ID: 2, Name: `2`}},
			httptest.NewRequest(http.MethodGet, `/?page=2&pageSize=1`, nil),
			`{"results":[{"id":2,"name":"2"}]}`,
			http.StatusOK,
		},
		{
			`should handle bad page`,
			nil,
			httptest.NewRequest(http.MethodGet, `/?page=invalid`, nil),
			`Error while parsing pagination: Error while parsing page param: strconv.ParseUint: parsing "invalid": invalid syntax
`,
			http.StatusBadRequest,
		},
		{
			`should handle bad page size`,
			nil,
			httptest.NewRequest(http.MethodGet, `/?pageSize=invalid`, nil),
			`Error while parsing pagination: Error while parsing pageSize param: strconv.ParseUint: parsing "invalid": invalid syntax
`,
			http.StatusBadRequest,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()
		users = testCase.init

		listCrud(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nlistCrud(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\nlistCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_getRequestID(t *testing.T) {
	var cases = []struct {
		intention string
		path      string
		want      uint
		wantErr   error
	}{
		{
			`should handle empty path`,
			`/`,
			0,
			fmt.Errorf(`strconv.ParseUint: parsing "": invalid syntax`),
		},
		{
			`should handle invalid number`,
			`/abc123`,
			0,
			fmt.Errorf(`strconv.ParseUint: parsing "abc123": invalid syntax`),
		},
		{
			`should handle positive number`,
			`/2`,
			2,
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		result, err := getRequestID(testCase.path)

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
			t.Errorf("%v\ngetRequestID(%v) = (%v, %v), want (%v, %v)", testCase.intention, testCase.path, result, err, testCase.want, testCase.wantErr)
		}
	}
}

func Test_readCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[uint]*user
		request    *http.Request
		id         uint
		want       string
		wantStatus int
	}{
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, nil),
			8000,
			``,
			http.StatusNotFound,
		},
		{
			`should return serialized instance`,
			map[uint]*user{1: {ID: 1, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			1,
			`{"id":1,"name":"test"}`,
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

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
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
		init       map[uint]*user
		request    *http.Request
		id         uint
		want       string
		wantStatus int
	}{
		{
			`should handle invalid JSON`,
			nil,
			httptest.NewRequest(http.MethodGet, `/1`, strings.NewReader(`{"name":"test"`)),
			1,
			`Error while unmarshalling body: unexpected end of JSON input
`,
			http.StatusBadRequest,
		},
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, strings.NewReader(`{"name":"Updated Test"}`)),
			8000,
			``,
			http.StatusNotFound,
		},
		{
			`should update given user`,
			map[uint]*user{1: {ID: 1, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, strings.NewReader(`{"name":"Updated Test"}`)),
			1,
			`{"id":1,"name":"Updated Test"}`,
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

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\nupdateCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_removeCrud(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[uint]*user
		request    *http.Request
		id         uint
		want       string
		wantStatus int
	}{
		{
			`should handle not found id`,
			nil,
			httptest.NewRequest(http.MethodGet, `/8000`, nil),
			8000,
			``,
			http.StatusNotFound,
		},
		{
			`should delete given user`,
			map[uint]*user{1: {ID: 1, Name: `test`}},
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			1,
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

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\nremoveCrudTest_removeCrud(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Test_ServeHTTP(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[uint]*user
		initSeq    uint
		request    *http.Request
		want       string
		wantStatus int
	}{
		{
			`should handle OPTIONS request for CORS`,
			nil,
			1,
			httptest.NewRequest(http.MethodOptions, `/`, nil),
			``,
			http.StatusNoContent,
		},
		{
			`should handle create request`,
			map[uint]*user{},
			1,
			httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(`{"name":"test"}`)),
			`{"id":1,"name":"test"}`,
			http.StatusCreated,
		},
		{
			`should handle list request`,
			map[uint]*user{1: {ID: 1, Name: `test`}},
			2,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`{"results":[{"id":1,"name":"test"}]}`,
			http.StatusOK,
		},
		{
			`should handle get request`,
			map[uint]*user{1: {ID: 1, Name: `test`}},
			2,
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			`{"id":1,"name":"test"}`,
			http.StatusOK,
		},
		{
			`should handle update request`,
			map[uint]*user{1: {ID: 1, Name: `test`}},
			2,
			httptest.NewRequest(http.MethodPut, `/1`, strings.NewReader(`{"name":"Updated test"}`)),
			`{"id":1,"name":"Updated test"}`,
			http.StatusOK,
		},
		{
			`should handle delete request`,
			map[uint]*user{1: {ID: 1, Name: `test`}},
			2,
			httptest.NewRequest(http.MethodDelete, `/1`, nil),
			``,
			http.StatusNoContent,
		},
		{
			`should handle unexpected method`,
			nil,
			1,
			httptest.NewRequest(http.MethodTrace, `/1`, nil),
			``,
			http.StatusMethodNotAllowed,
		},
	}

	for _, testCase := range cases {
		writer := httptest.NewRecorder()

		users = testCase.init
		seq = testCase.initSeq

		Handler().ServeHTTP(writer, testCase.request)

		if result := writer.Code; result != testCase.wantStatus {
			t.Errorf("%v\nServeHTTP(%v) = %v, want status %v", testCase.intention, testCase.request, result, testCase.wantStatus)
		}

		if result, _ := httputils.ReadBody(writer.Result().Body); string(result) != testCase.want {
			t.Errorf("%v\nServeHTTP(%v) = %v, want %v", testCase.intention, testCase.request, string(result), testCase.want)
		}
	}
}

func Benchmark_ServeHTTP_list(b *testing.B) {
	handler := Handler()
	users = map[uint]*user{}
	seq = 1

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, `/`, nil))
	}
}

func Benchmark_ServeHTTP_options(b *testing.B) {
	handler := Handler()
	users = map[uint]*user{}
	seq = 1

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodOptions, `/`, nil))
	}
}

func Benchmark_ServeHTTP_create(b *testing.B) {
	handler := Handler()
	users = map[uint]*user{}
	seq = 1

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(`{"name":"test"}`)))
	}
}

func Benchmark_ServeHTTP_get(b *testing.B) {
	handler := Handler()
	users = map[uint]*user{1: {ID: 1, Name: `test`}}
	seq = 2

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, `/1`, nil))
	}
}
