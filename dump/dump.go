package dump

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils/request"
)

// Handler for dump request. Should be use with net/http
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headers string
		for key, value := range r.Header {
			headers += fmt.Sprintf("%s: %s\n", key, strings.Join(value, `,`))
		}

		var params string
		for key, value := range r.URL.Query() {
			params += fmt.Sprintf("%s: %s\n", key, strings.Join(value, `,`))
		}

		var form string
		if err := r.ParseForm(); err != nil {
			form = fmt.Sprintf(`Error while parsing form: %v`, err)
		} else {
			for key, value := range r.PostForm {
				form += fmt.Sprintf("%s: %s\n", key, strings.Join(value, `,`))
			}
		}

		body, err := request.ReadBody(r.Body)
		if err != nil {
			log.Printf(`Error while reading body: %v`, err)
		}

		outputPattern := "\n%s %s\n"
		outputData := []interface{}{
			r.Method,
			r.URL.Path,
		}

		if headers != `` {
			outputPattern += "Headers\n%s\n"
			outputData = append(outputData, headers)
		}

		if params != `` {
			outputPattern += "Params\n%s\n"
			outputData = append(outputData, params)
		}

		if form != `` {
			outputPattern += "Form\n%s\n"
			outputData = append(outputData, form)
		}

		if len(body) != 0 {
			outputPattern += "Body\n%s\n"
			outputData = append(outputData, body)
		}

		log.Printf(outputPattern, outputData...)
	})
}
