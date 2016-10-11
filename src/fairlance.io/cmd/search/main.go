package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	respond "gopkg.in/matryer/respond.v1"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

var (
	jobsIndex        bleve.Index
	freelancersIndex bleve.Index
	indicesDir       = *flag.String("indicesDir", "/tmp/indices", "Location where the indices are located.")
	port             = *flag.String("port", "3002", "Port.")
	respondOptions   *respond.Options
)

func init() {
	flag.Parse()

	var err error
	jobsIndex, err = getIndex("jobs")
	if err != nil {
		log.Fatal(err)
	}

	freelancersIndex, err = getIndex("freelancers")
	if err != nil {
		log.Fatal(err)
	}

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
	searchRequest, err := getSearchRequest(r)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

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
	searchRequest.Fields = []string{"tags"}

	tagsFacet := bleve.NewFacetRequest("tags", 99999)
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

func getSearchRequest(r *http.Request) (*bleve.SearchRequest, error) {
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
