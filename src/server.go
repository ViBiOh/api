package main

import "net/http"
import "html"
import "encoding/json"
import "time"
import "log"
import "strings"
import "strconv"
import "regexp"
import "io/ioutil"

const delayInSeconds = 1
const port = "1080"

const PERFORMANCE_URL = "http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id="
var PERF_ONE_MONTH = regexp.MustCompile("<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>")
var PERF_THREE_MONTH = regexp.MustCompile("<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>")
var PERF_SIX_MONTH = regexp.MustCompile("<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>")
var PERF_ONE_YEAR = regexp.MustCompile("<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>")

type Hello struct {
	Name string `json:"greeting"`
}

func pluralize(s string, n int) string {
	if n > 1 {
		return (s + "s")
	}
	return s
}

func apiHello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(delayInSeconds * time.Second)
	hello := Hello{"Hello " + html.EscapeString(strings.Replace(r.URL.Path, "/hello/", "", -1)) + ", I'm greeting you from the server with " + strconv.Itoa(delayInSeconds) + " " + pluralize("second", delayInSeconds) + " delay"}

	helloJson, err := json.Marshal(hello)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(helloJson)
	} else {
		log.Fatal(err)
	}
}

type Performance struct {
	OneMonth string
}

func apiPerf(w http.ResponseWriter, r *http.Request) {
  morningStarId := strings.ToLower(strings.Replace(r.URL.Path, "/perf/", "", -1))
	response, err := http.Get(PERFORMANCE_URL + morningStarId)

	if err != nil {
		http.Error(w, "Error while fetching data", 500)
		return
	}

  if response.StatusCode >= 400 {
    http.Error(w, morningStarId + " not found", response.StatusCode)
    return
  }

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error while reading body", 500)
		return
	}
	
	performance := Performance{string(PERF_ONE_MONTH.FindSubmatch(body)[1][:])}
	performanceJson, errJson := json.Marshal(performance)
	
	if errJson != nil {
		http.Error(w, "Error while marshalling json", 500)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(performanceJson)
	}
}

func main() {
	http.HandleFunc("/hello/", apiHello)
	http.HandleFunc("/perf/", apiPerf)

	log.Print("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
