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

	for _, tag := range r.URL.Query()["tags"] {
		shoulds = append(shoulds, bleve.NewMatchQuery(tag))
	}

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
