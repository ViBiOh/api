package crud

import (
	"net/http"
	"testing"
)

func Test_getID(t *testing.T) {
	emptyRequest, _ := http.NewRequest(http.MethodGet, `/`, nil)
	simpleRequest, _ := http.NewRequest(http.MethodGet, `/abc-1234`, nil)
	complexRequest, _ := http.NewRequest(http.MethodGet, `/def-5678/links/`, nil)

	var cases = []struct {
		intention string
		request   *http.Request
		want      string
	}{
		{
			`should work with empty URL`,
			emptyRequest,
			``,
		},
		{
			`should work with simple ID URL`,
			simpleRequest,
			`abc-1234`,
		},
		{
			`should work with complex ID URL`,
			complexRequest,
			`def-5678`,
		},
	}

	for _, testCase := range cases {
		if result := getID(testCase.request); result != testCase.want {
			t.Errorf("%s\ngetID(`%+v`) = %+v, want %+v", testCase.intention, testCase.request, result, testCase.want)
		}
	}
}
