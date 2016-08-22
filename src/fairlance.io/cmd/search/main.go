package main

import (
	"fmt"
	"net/http"

	respond "gopkg.in/matryer/respond.v1"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
)

var (
	jobsIndex        bleve.Index
	freelancersIndex bleve.Index
)

func init() {
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
	http.HandleFunc("/jobs", Jobs)
	http.HandleFunc("/freelancers", Freelancers)

	fmt.Println("Listening on port 3002")
	panic(http.ListenAndServe(":3002", nil))
}

func Jobs(w http.ResponseWriter, r *http.Request) {
	searchRequest := getSearchRequest()

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

func Freelancers(w http.ResponseWriter, r *http.Request) {
	searchRequest := getSearchRequest()

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
	fmt.Printf("Opening %sIndex ...\n", dbName)
	index, err := bleve.Open("/tmp/" + dbName)
	if err != nil {
		return index, err
	}

	fmt.Printf("Opened %sIndex\n", dbName)
	return index, nil
}

func getSearchRequest() *bleve.SearchRequest {
	// search for some text
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	// searchRequest.Highlight = bleve.NewHighlight()

	return searchRequest
}
