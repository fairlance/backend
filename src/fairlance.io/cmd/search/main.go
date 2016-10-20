package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	respond "gopkg.in/matryer/respond.v1"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

var (
	port           = *flag.String("port", "3002", "Port.")
	searcherURL    = *flag.String("searcherURL", "http://localhost:3003", "Url of the searcher.")
	respondOptions *respond.Options
)

func init() {
	flag.Parse()

	f, err := os.OpenFile("/var/log/fairlance/search.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)

	respondOptions = &respond.Options{
		Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
			dataEnvelope := map[string]interface{}{"code": status}
			if err, ok := data.(error); ok {
				dataEnvelope["error"] = err.Error()
				dataEnvelope["success"] = false
			} else {
				dataEnvelope["data"] = data
				dataEnvelope["success"] = true
			}
			return status, dataEnvelope
		},
	}
}

func main() {
	http.Handle("/job", corsHandler(respondOptions.Handler(http.HandlerFunc(jobs))))
	http.Handle("/job/tags", corsHandler(respondOptions.Handler(http.HandlerFunc(jobTags))))
	http.Handle("/freelancer", corsHandler(respondOptions.Handler(http.HandlerFunc(freelancers))))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func jobs(w http.ResponseWriter, r *http.Request) {
	searchRequest, err := getJobSearchRequest(r)
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	jobsSearchResults, err := doRequest("jobs", searchRequest)
	if err != nil {
		respond.With(w, r, http.StatusBadGateway, err)
		return
	}

	jobs := []interface{}{}
	for _, hit := range jobsSearchResults.Hits {
		jobs = append(jobs, hit.Fields)
	}

	respond.With(w, r, http.StatusOK, struct {
		Total int           `json:"total"`
		Items []interface{} `json:"items"`
	}{
		Total: len(jobs),
		Items: jobs,
	})
}

func jobTags(w http.ResponseWriter, r *http.Request) {

	query := bleve.NewMatchAllQuery()

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"tags"}

	tagsFacet := bleve.NewFacetRequest("tags", 99999)
	searchRequest.AddFacet("tags", tagsFacet)

	jobsSearchResults, err := doRequest("jobs", searchRequest)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	tags := []string{}
	for _, t := range jobsSearchResults.Facets["tags"].Terms {
		tags = append(tags, fmt.Sprintf("%s", t.Term))
	}

	respond.With(w, r, http.StatusOK, struct {
		Total int      `json:"total"`
		Tags  []string `json:"tags"`
	}{
		Total: len(tags),
		Tags:  tags,
	})
}

func freelancers(w http.ResponseWriter, r *http.Request) {

	searchRequest := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
	searchRequest.Fields = []string{"*"}

	freelnacersSearchResults, err := doRequest("freelancers", searchRequest)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	freelancers := []interface{}{}
	for _, hit := range freelnacersSearchResults.Hits {
		freelancers = append(freelancers, hit.Fields)
	}
	respond.With(w, r, http.StatusOK, struct {
		Total int           `json:"total"`
		Items []interface{} `json:"items"`
	}{
		Total: len(freelancers),
		Items: freelancers,
	})
}

func getJobSearchRequest(r *http.Request) (*bleve.SearchRequest, error) {
	musts := []query.Query{}
	// mustNots := []query.Query{}
	// shoulds := []query.Query{}

	tagShoulds := []query.Query{}
	for _, tag := range r.URL.Query()["tags"] {
		tagShoulds = append(tagShoulds, bleve.NewMatchQuery(tag))
	}
	if len(tagShoulds) > 0 {
		booleanQuery := bleve.NewBooleanQuery()
		booleanQuery.AddShould(tagShoulds...)
		musts = append(musts, booleanQuery)
	}

	value1 := 0.0
	if len(r.URL.Query().Get("price_from")) != 0 {
		intValue1, err := strconv.ParseInt(r.URL.Query().Get("price_from"), 10, 64)
		if err != nil {
			return nil, err
		}
		value1 = float64(intValue1)
	}
	value2 := math.MaxFloat64
	if len(r.URL.Query().Get("price_to")) != 0 {
		intValue2, err := strconv.ParseInt(r.URL.Query().Get("price_to"), 10, 64)
		if err != nil {
			return nil, err
		}
		value2 = float64(intValue2)
	}

	inclusiveValue1 := true
	inclusiveValue2 := false
	numericRangeIncludiveQuery := bleve.NewNumericRangeInclusiveQuery(
		&value1,
		&value2,
		&inclusiveValue1,
		&inclusiveValue2,
	)
	numericRangeIncludiveQuery.SetField("price")
	musts = append(musts, numericRangeIncludiveQuery)

	period := int64(30)
	if len(r.URL.Query().Get("period")) != 0 {
		periodTemp, err := strconv.ParseInt(r.URL.Query().Get("period"), 10, 64)
		if err != nil {
			return nil, err
		}
		if periodTemp > 0 && periodTemp <= 365 {
			period = periodTemp
		}
	}

	now := time.Now()
	dateTo := time.Now().Add(time.Duration(24*period) * time.Hour)
	dateRangeQuery := bleve.NewDateRangeQuery(now, dateTo)
	dateRangeQuery.SetField("startDate")
	musts = append(musts, dateRangeQuery)

	query := bleve.NewBooleanQuery()
	query.AddMust(musts...)
	// query.AddMustNot(mustNots...)
	// query.AddShould(shoulds...)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}

	return searchRequest, nil
}

func doRequest(index string, searchRequest *bleve.SearchRequest) (bleve.SearchResult, error) {
	var jobsSearchResults bleve.SearchResult
	jsonBytes, err := json.Marshal(searchRequest)
	if err != nil {
		return jobsSearchResults, err
	}

	url := searcherURL + "/api/" + index + "/_search"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return jobsSearchResults, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return jobsSearchResults, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jobsSearchResults, err
	}

	if resp.StatusCode != http.StatusOK {
		return jobsSearchResults, errors.New(string(body))
	}

	err = json.Unmarshal(body, &jobsSearchResults)
	if err != nil {
		return jobsSearchResults, err
	}

	return jobsSearchResults, nil
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			// todo: make configurable
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		// Stop here for a Preflighted OPTIONS request.
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
