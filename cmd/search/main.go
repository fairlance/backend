package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	respond "gopkg.in/matryer/respond.v1"

	"github.com/fairlance/backend/middleware"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

func main() {
	log.SetFlags(log.Lshortfile)
	var port = os.Getenv("PORT")
	var searcherURL = os.Getenv("SEARCHER_URL")
	corsJSON := middleware.Chain(
		middleware.CORSHandler,
		middleware.JSONEnvelope,
	)
	http.Handle("/job", corsJSON(jobs(searcherURL)))
	http.Handle("/job/tags", corsJSON(jobTags(searcherURL)))
	http.Handle("/freelancer", corsJSON(freelancers(searcherURL)))
	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func jobs(searcherURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		searchRequest, err := getJobSearchRequest(r)
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}

		jobsSearchResults, err := doRequest(searcherURL, "jobs", searchRequest)
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
	})
}

func jobTags(searcherURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		query := bleve.NewMatchAllQuery()

		searchRequest := bleve.NewSearchRequest(query)
		searchRequest.Fields = []string{"tags"}

		tagsFacet := bleve.NewFacetRequest("tags", 99999)
		searchRequest.AddFacet("tags", tagsFacet)

		jobsSearchResults, err := doRequest(searcherURL, "jobs", searchRequest)
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
	})
}

func freelancers(searcherURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		searchRequest := bleve.NewSearchRequest(bleve.NewMatchAllQuery())
		searchRequest.Fields = []string{"*"}

		freelnacersSearchResults, err := doRequest(searcherURL, "freelancers", searchRequest)
		if err != nil {
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		freelancers := []interface{}{}
		for _, hit := range freelnacersSearchResults.Hits {
			freelancers = append(freelancers, hit.Fields)
		}
		respond.With(w, r, http.StatusOK, struct {
			Total uint64        `json:"total"`
			Items []interface{} `json:"items"`
		}{
			Total: freelnacersSearchResults.Total,
			Items: freelancers,
		})
	})
}

func getJobSearchRequest(r *http.Request) (*bleve.SearchRequest, error) {
	size := 50
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

	var priceFromVal *float64
	var priceToVal *float64
	if r.URL.Query().Get("price_from") != "" {
		priceFromFloat, err := strconv.ParseFloat(r.URL.Query().Get("price_from"), 64)
		if err != nil {
			return nil, err
		}
		priceFromVal = &priceFromFloat
	}

	if r.URL.Query().Get("price_to") != "" {
		priceToFloat, err := strconv.ParseFloat(r.URL.Query().Get("price_to"), 64)
		if err != nil {
			return nil, err
		}
		priceToVal = &priceToFloat
	}

	if priceFromVal != nil || priceToVal != nil {
		inclusiveValue1 := true
		inclusiveValue2 := false
		numericRangeIncludiveQuery := bleve.NewNumericRangeInclusiveQuery(
			priceFromVal,
			priceToVal,
			&inclusiveValue1,
			&inclusiveValue2,
		)
		numericRangeIncludiveQuery.SetField("price")
		musts = append(musts, numericRangeIncludiveQuery)
	}

	if len(r.URL.Query().Get("period")) != 0 {
		period, err := strconv.Atoi(r.URL.Query().Get("period"))
		if err != nil {
			return nil, err
		}

		if period < 0 || period > 365 {
			period = 30
		}
		now := time.Now()
		dateTo := time.Now().Add(time.Duration(24*period) * time.Hour)
		dateRangeQuery := bleve.NewDateRangeQuery(now, dateTo)
		dateRangeQuery.SetField("startDate")
		musts = append(musts, dateRangeQuery)
	}

	var query query.Query
	if len(musts) > 0 {
		query := bleve.NewBooleanQuery()
		query.AddMust(musts...)
	} else {
		query = bleve.NewMatchAllQuery()
	}
	// query.AddMustNot(mustNots...)
	// query.AddShould(shoulds...)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.SortBy([]string{"-updatedAt"})
	searchRequest.Size = size
	if len(r.URL.Query().Get("page")) != 0 {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			return nil, err
		}
		if page > 0 {
			searchRequest.From = (page - 1) * size
		}
	}
	searchRequest.Fields = []string{"*"}

	return searchRequest, nil
}

func doRequest(searcherURL, index string, searchRequest *bleve.SearchRequest) (bleve.SearchResult, error) {
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return jobsSearchResults, err
	}
	defer resp.Body.Close()

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
