package morningStar

import "net/http"
import "strings"
import "strconv"
import "regexp"
import "io/ioutil"
import "../json"

const PERFORMANCE_URL = "http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=1&id="
const VOLATILITE_URL = "http://www.morningstar.fr/fr/funds/snapshot/snapshot.aspx?tab=2&id="
const SEARCH_ID = "http://www.morningstar.fr/fr/util/SecuritySearch.ashx?q="

var PERF_ONE_MONTH = regexp.MustCompile("<td[^>]*?>1 mois</td><td[^>]*?>(.*?)</td>")
var PERF_THREE_MONTH = regexp.MustCompile("<td[^>]*?>3 mois</td><td[^>]*?>(.*?)</td>")
var PERF_SIX_MONTH = regexp.MustCompile("<td[^>]*?>6 mois</td><td[^>]*?>(.*?)</td>")
var PERF_ONE_YEAR = regexp.MustCompile("<td[^>]*?>1 an</td><td[^>]*?>(.*?)</td>")
var VOL_3_YEAR = regexp.MustCompile("<td[^>]*?>Ecart-type 3 ans.?</td><td[^>]*?>(.*?)</td>")

type Performance struct {
  MorningStarId string  `json:"id"`
  OneMonth      float64 `json:"1m"`
  ThreeMonth    float64 `json:"3m"`
  SixMonth      float64 `json:"6m"`
  OneYear       float64 `json:"1y"`
  VolThreeYears       float64 `json:"v1y"`
}

func getPerformance(extract *regexp.Regexp, body []byte) float64 {
  match := extract.FindSubmatch(body)
  if (match == nil) {
    return 0.0
  }

  rawResult := string(match[1][:])
  dotResult := strings.Replace(rawResult, ",", ".", -1)
  percentageResult := strings.Replace(dotResult, "%", "", -1)
  trimResult := strings.TrimSpace(percentageResult)

  result, err := strconv.ParseFloat(trimResult, 64)
  if err != nil {
    return 0.0
  }
  return result
}

func Handler(w http.ResponseWriter, r *http.Request) {
  morningStarId := strings.ToLower(strings.Replace(r.URL.Path, "/morningStar/", "", -1))

  if (strings.TrimSpace(morningStarId) == ""){
    http.Error(w, "Missing id", 400)
    return
  } 

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

  oneMonth := getPerformance(PERF_ONE_MONTH, performanceBody)
  threeMonths := getPerformance(PERF_THREE_MONTH, performanceBody)
  sixMonths := getPerformance(PERF_SIX_MONTH, performanceBody)
  oneYear := getPerformance(PERF_ONE_YEAR, performanceBody)
  volThreeYears := getPerformance(VOL_3_YEAR, volatiliteBody)

  performance := Performance{morningStarId, oneMonth, threeMonths, sixMonths, oneYear, volThreeYears}
  json.ResponseJson(w, performance)
}