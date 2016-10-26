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
const VOLATILITE_URL = "http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id="
const SEARCH_ID = "http://www.morningstar.fr/fr/util/SecuritySearch.ashx?q="

var PERF_ONE_MONTH = regexp.MustCompile("<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>")
var PERF_THREE_MONTH = regexp.MustCompile("<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>")
var PERF_SIX_MONTH = regexp.MustCompile("<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>")
var PERF_ONE_YEAR = regexp.MustCompile("<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>")
var VOL_3_YEAR = regexp.MustCompile("<td[^>]*?>Ecart-type 3 ans.?</td><td[^>]*?>(.*?)</td>")

func responseJson(w http.ResponseWriter, obj interface{}) {
	objJson, err := json.Marshal(obj)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(objJson)
	} else {
		http.Error(w, "Error while marshalling JSON response", 500)
	}
}

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

	responseJson(w, hello)
}

type Performance struct {
	MorningStarId string  `json:"id"`
	OneMonth      float64 `json:"1m"`
	ThreeMonth    float64 `json:"3m"`
	SixMonth      float64 `json:"6m"`
	OneYear       float64 `json:"1y"`
	VolThreeYears       float64 `json:"v1y"`
}

func getPerformance(rawValue []byte) float64 {
	result, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(strings.Replace(string(rawValue[:]), ",", ".", -1), "%", "", -1)), 64)
	return result
}

func apiPerf(w http.ResponseWriter, r *http.Request) {
	morningStarId := strings.ToLower(strings.Replace(r.URL.Path, "/perf/", "", -1))

	performanceResponse, err := http.Get(PERFORMANCE_URL + morningStarId)
	if err != nil {
		http.Error(w, "Error while fetching data", 500)
		return
	}

	if performanceResponse.StatusCode >= 400 {
		http.Error(w, morningStarId+" not found", performanceResponse.StatusCode)
		return
	}

	defer performanceResponse.Body.Close()
	performanceBody, err := ioutil.ReadAll(performanceResponse.Body)
	if err != nil {
		http.Error(w, "Error while reading performance body", 500)
		return
	}
	
	volatiliteResponse, err := http.Get(VOLATILITE_URL + morningStarId)
	if err != nil {
		http.Error(w, "Error while fetching data", 500)
		return
	}

	defer volatiliteResponse.Body.Close()
	volatiliteBody, err := ioutil.ReadAll(volatiliteResponse.Body)
	if err != nil {
		http.Error(w, "Error while reading volatilite body", 500)
		return
	}

	oneMonth := getPerformance(PERF_ONE_MONTH.FindSubmatch(performanceBody)[1])
	threeMonths := getPerformance(PERF_THREE_MONTH.FindSubmatch(performanceBody)[1])
	sixMonths := getPerformance(PERF_SIX_MONTH.FindSubmatch(performanceBody)[1])
	oneYear := getPerformance(PERF_ONE_YEAR.FindSubmatch(performanceBody)[1])
	volThreeYears := getPerformance(VOL_3_YEAR.FindSubmatch(volatiliteBody)[1])

	performance := Performance{morningStarId, oneMonth, threeMonths, sixMonths, oneYear, volThreeYears}
	responseJson(w, performance)
}

func main() {
	http.HandleFunc("/hello/", apiHello)
	http.HandleFunc("/perf/", apiPerf)

	log.Print("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
