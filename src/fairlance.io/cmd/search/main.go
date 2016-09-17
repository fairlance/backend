package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	respond "gopkg.in/matryer/respond.v1"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
)

var (
	jobsIndex        bleve.Index
	freelancersIndex bleve.Index
	indicesDir       = *flag.String("indicesDir", "/tmp", "Location where the indices are located.")
	port             = *flag.String("port", "3002", "Port.")
)

func init() {
	flag.Parse()

	var err error
	jobsIndex, err = getIndex("jobs")
	if err != nil {
		panic(err)
	}

	freelancersIndex, err = getIndex("freelancers")
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/jobs", CORSHandler(http.HandlerFunc(jobs)))
	http.Handle("/jobs/tags", CORSHandler(http.HandlerFunc(jobTags)))
	http.Handle("/freelancers", CORSHandler(http.HandlerFunc(freelancers)))

	panic(http.ListenAndServe(":"+port, nil))
}

func jobs(w http.ResponseWriter, r *http.Request) {
	searchRequest := getSearchRequest(r)

	jobsSearchResults, err := jobsIndex.Search(searchRequest)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
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
	searchRequest.Fields = []string{"tags.name"}

	tagsFacet := bleve.NewFacetRequest("tags.name", 99999)
	searchRequest.AddFacet("tags", tagsFacet)

	jobsSearchResults, err := jobsIndex.Search(searchRequest)
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

	freelnacersSearchResults, err := freelancersIndex.Search(searchRequest)
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

func getIndex(dbName string) (bleve.Index, error) {
	index, err := bleve.Open(indicesDir + "/" + dbName)
	if err != nil {
		return index, err
	}

	return index, nil
}

func getSearchRequest(r *http.Request) *bleve.SearchRequest {
	musts := []bleve.Query{}
	mustNots := []bleve.Query{}
	shoulds := []bleve.Query{}

	tagShoulds := []bleve.Query{}
	for _, tag := range r.URL.Query()["tags"] {
		tagShoulds = append(tagShoulds, bleve.NewMatchQuery(tag))
	}
	if len(tagShoulds) > 0 {
		musts = append(musts, bleve.NewBooleanQuery(nil, tagShoulds, nil))
	}

	value1 := 0.0
	if len(r.URL.Query().Get("price_from")) != 0 {
		intValue1, err := strconv.ParseInt(r.URL.Query().Get("price_from"), 10, 64)
		if err != nil {
			panic(err)
		}
		value1 = float64(intValue1)
	}
	value2 := math.MaxFloat64
	if len(r.URL.Query().Get("price_to")) != 0 {
		intValue2, err := strconv.ParseInt(r.URL.Query().Get("price_to"), 10, 64)
		if err != nil {
			panic(err)
		}
		value2 = float64(intValue2)
	}

	inclusiveValue1 := true
	inclusiveValue2 := false
	musts = append(musts, bleve.NewNumericRangeInclusiveQuery(&value1, &value2, &inclusiveValue1, &inclusiveValue2).SetField("price"))

	period := int64(30)
	if len(r.URL.Query().Get("period")) != 0 {
		periodTemp, err := strconv.ParseInt(r.URL.Query().Get("period"), 10, 64)
		if err != nil {
			panic(err)
		}
		if periodTemp > 0 && periodTemp <= 365 {
			period = periodTemp
		}
	}

	now := time.Now().Format(time.RFC3339)
	dateTo := time.Now().Add(time.Duration(24*period) * time.Hour).Format(time.RFC3339)
	musts = append(musts, bleve.NewDateRangeQuery(&now, &dateTo).SetField("startDate"))

	query := bleve.NewBooleanQuery(musts, shoulds, mustNots)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}

	return searchRequest
}

// CORSHandler handler
func CORSHandler(next http.Handler) http.Handler {
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
