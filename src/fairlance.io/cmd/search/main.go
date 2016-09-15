package main

import (
	"flag"
	"fmt"
	"net/http"

	respond "gopkg.in/matryer/respond.v1"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
)

var (
	jobsIndex        bleve.Index
	freelancersIndex bleve.Index
	indicesDir       = flag.String("indicesDir", "/tmp", "Location where the indices are located.")
	port             = flag.String("port", "3002", "Port.")
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
	http.HandleFunc("/jobs", jobs)
	http.HandleFunc("/freelancers", freelancers)

	panic(http.ListenAndServe(":"+*port, nil))
}

func jobs(w http.ResponseWriter, r *http.Request) {
	searchRequest := getSearchRequest(r)

	tagsFacet := bleve.NewFacetRequest("tags.name", 100)
	searchRequest.AddFacet("tags", tagsFacet)

	jobsSearchResults, err := jobsIndex.Search(searchRequest)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	jobs := []interface{}{}
	for _, hit := range jobsSearchResults.Hits {
		jobs = append(jobs, hit.Fields)
	}

	tags := []string{}
	for _, t := range jobsSearchResults.Facets["tags"].Terms {
		tags = append(tags, fmt.Sprintf("%s (%d)", t.Term, t.Count))
	}

	respond.With(w, r, http.StatusOK, struct {
		Total int           `json:"total"`
		Items []interface{} `json:"items"`
		Tags  []string      `json:"tags"`
	}{
		Total: len(jobs),
		Items: jobs,
		Tags:  tags,
	})
}

func freelancers(w http.ResponseWriter, r *http.Request) {
	searchRequest := getSearchRequest(r)

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
	index, err := bleve.Open(*indicesDir + "/" + dbName)
	if err != nil {
		return index, err
	}

	return index, nil
}

func getSearchRequest(r *http.Request) *bleve.SearchRequest {
	musts := []bleve.Query{}
	mustNots := []bleve.Query{}
	shoulds := []bleve.Query{}

	for _, tag := range r.URL.Query()["tags"] {
		shoulds = append(shoulds, bleve.NewMatchQuery(tag))
	}

	query := bleve.NewBooleanQuery(musts, shoulds, mustNots)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}

	return searchRequest
}
