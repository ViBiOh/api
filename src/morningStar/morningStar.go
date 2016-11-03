package morningStar

import "net/http"
import "errors"
import "time"
import "strings"
import "strconv"
import "regexp"
import "io/ioutil"
import "encoding/json"
import "../jsonHttp"

const PERFORMANCE_URL = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id=`
const VOLATILITE_URL = `http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id=`
const SEARCH_ID = `http://www.morningstar.fr/fr/util/SecuritySearch.ashx?q=`

var ISIN_REQUEST = regexp.MustCompile(`(.+?)/isin`)
var PERF_REQUEST = regexp.MustCompile(`(.+?)`)
var LIST_REQUEST = regexp.MustCompile(`^/?$`)
var CARRIAGE_RETURN = regexp.MustCompile(`\r?\n`)
var END_CARRIAGE_RETURN = regexp.MustCompile(`\r?\n$`)
var PIPE = regexp.MustCompile(`[|]`)

var PERF_ONE_MONTH = regexp.MustCompile(`<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_THREE_MONTH = regexp.MustCompile(`<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_SIX_MONTH = regexp.MustCompile(`<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>`)
var PERF_ONE_YEAR = regexp.MustCompile(`<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>`)
var VOL_3_YEAR = regexp.MustCompile(`<td[^>]*?>Ecart-type 3 ans.?</td><td[^>]*?>(.*?)</td>`)

var PERFORMANCE_CACHE = make(map[string]Performance)

type Performance struct {
	MorningStarId string    `json:"id"`
	OneMonth      float64   `json:"1m"`
	ThreeMonth    float64   `json:"3m"`
	SixMonth      float64   `json:"6m"`
	OneYear       float64   `json:"1y"`
	VolThreeYears float64   `json:"v1y"`
	Update        time.Time `json:"ts"`
}

type Search struct {
	Id    string `json:"i"`
	Label string `json:"n"`
}

type Results struct {
	Results interface{} `json:"results"`
}

func getBody(url string, w http.ResponseWriter) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New(`Error while retrieving data from `+url)
	}

	if response.StatusCode >= 400 {
		return nil, errors.New(`Got error `+strconv.Itoa(response.StatusCode)+` while getting `+url)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(`Error while reading body of `+url)
	}

	return body, nil
}

func getPerformance(extract *regexp.Regexp, body []byte) float64 {
	match := extract.FindSubmatch(body)
	if match == nil {
		return 0.0
	}

	rawResult := string(match[1][:])
	dotResult := strings.Replace(rawResult, `,`, `.`, -1)
	percentageResult := strings.Replace(dotResult, `%`, ``, -1)
	trimResult := strings.TrimSpace(percentageResult)

	result, err := strconv.ParseFloat(trimResult, 64)
	if err != nil {
		return 0.0
	}
	return result
}

func singlePerformance(w http.ResponseWriter, morningStarId string) (*Performance, error) {
	performance, present := PERFORMANCE_CACHE[morningStarId]
	if present && time.Now().Add(time.Hour*-18).Before(performance.Update) {
		return &performance, nil
	}

	performanceBody, err := getBody(PERFORMANCE_URL+morningStarId, w)
	if err != nil {
		return nil, err
	}

	volatiliteBody, err := getBody(VOLATILITE_URL+morningStarId, w)
	if err != nil {
		return nil, err
	}

	oneMonth := getPerformance(PERF_ONE_MONTH, performanceBody)
	threeMonths := getPerformance(PERF_THREE_MONTH, performanceBody)
	sixMonths := getPerformance(PERF_SIX_MONTH, performanceBody)
	oneYear := getPerformance(PERF_ONE_YEAR, performanceBody)
	volThreeYears := getPerformance(VOL_3_YEAR, volatiliteBody)

	performance = Performance{morningStarId, oneMonth, threeMonths, sixMonths, oneYear, volThreeYears, time.Now()}
	PERFORMANCE_CACHE[morningStarId] = performance

	return &performance, nil
}

func singlePerformanceHandler(w http.ResponseWriter, morningStarId string) {
	performance, err := singlePerformance(w, morningStarId)

	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		jsonHttp.ResponseJson(w, *performance)
	}
}

func isinHandler(w http.ResponseWriter, isin string) {
	searchBody, err := getBody(SEARCH_ID+strings.ToLower(isin), w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	cleanBody := END_CARRIAGE_RETURN.ReplaceAllString(string(searchBody[:]), ``)
	lines := CARRIAGE_RETURN.Split(cleanBody, -1)
	size := len(lines)

	results := make([]Search, size)
	for i := 0; i < size; i++ {
		jsonErr := json.Unmarshal([]byte(PIPE.Split(lines[i], -1)[1]), &results[i])
		if jsonErr != nil {
			http.Error(w, `Error while unmarshalling data for ISIN `+isin, 500)
		}
	}

	jsonHttp.ResponseJson(w, Results{results})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.ToLower(strings.Replace(r.URL.Path, `/morningStar/`, ``, -1))

	if PERF_REQUEST.MatchString(path) {
			singlePerformanceHandler(w, path)
	} else if ISIN_REQUEST.MatchString(path) {
		isinHandler(w, ISIN_REQUEST.FindStringSubmatch(path)[1])
	}
}
