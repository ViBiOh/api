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
		init       map[int64]*user
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
			map[int64]*user{1: {ID: 1, Name: `1`}, 2: {ID: 2, Name: `2`}},
			httptest.NewRequest(http.MethodGet, `/?page=2&pageSize=1`, nil),
			`{"results":[{"id":2,"name":"2"}]}`,
			http.StatusOK,
		},
		{
			`should handle bad page`,
			nil,
			httptest.NewRequest(http.MethodGet, `/?page=invalid`, nil),
			`Error while parsing page param: strconv.ParseInt: parsing "invalid": invalid syntax
`,
			http.StatusBadRequest,
		},
		{
			`should handle bad page`,
			nil,
			httptest.NewRequest(http.MethodGet, `/?pageSize=invalid`, nil),
			`Error while parsing pageSize param: strconv.ParseInt: parsing "invalid": invalid syntax
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
			httptest.NewRequest(http.MethodGet, `/1`, strings.NewReader(`{"name":"test"`)),
			`Error while unmarshalling body: unexpected end of JSON input
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

func Test_ServeHTTP(t *testing.T) {
	var cases = []struct {
		intention  string
		init       map[int64]*user
		initSeq    int64
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
			map[int64]*user{},
			1,
			httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(`{"name":"test"}`)),
			`{"id":1,"name":"test"}`,
			http.StatusCreated,
		},
		{
			`should handle list request`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
			2,
			httptest.NewRequest(http.MethodGet, `/`, nil),
			`{"results":[{"id":1,"name":"test"}]}`,
			http.StatusOK,
		},
		{
			`should handle get request`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
			2,
			httptest.NewRequest(http.MethodGet, `/1`, nil),
			`{"id":1,"name":"test"}`,
			http.StatusOK,
		},
		{
			`should handle update request`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
			2,
			httptest.NewRequest(http.MethodPut, `/1`, strings.NewReader(`{"name":"Updated test"}`)),
			`{"id":1,"name":"Updated test"}`,
			http.StatusOK,
		},
		{
			`should handle delete request`,
			map[int64]*user{1: {ID: 1, Name: `test`}},
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
	users = map[int64]*user{}
	seq = 1

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, `/`, nil))
	}
}

func Benchmark_ServeHTTP_options(b *testing.B) {
	handler := Handler()
	users = map[int64]*user{}
	seq = 1

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodOptions, `/`, nil))
	}
}

func Benchmark_ServeHTTP_create(b *testing.B) {
	handler := Handler()
	users = map[int64]*user{}
	seq = 1

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(`{"name":"test"}`)))
	}
}

func Benchmark_ServeHTTP_get(b *testing.B) {
	handler := Handler()
	users = map[int64]*user{1: {ID: 1, Name: `test`}}
	seq = 2

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, `/1`, nil))
	}
}
