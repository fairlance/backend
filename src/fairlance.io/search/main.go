package main

import (
	"encoding/json"
	"flag"
	"gopkg.in/matryer/respond.v1"
	"io/ioutil"
	"net/http"
	"strconv"
)

var port int
var elasticsearchUrl string

type ElasticSearchResponse struct {
	Hits Hits
}
type Hits struct {
	Hits  []Source
	Total int
}
type Source struct {
	Source map[string]interface{} `json:"_source"`
}

type processFunc func(map[string]interface{}) interface{}

func main() {
	flag.IntVar(&port, "port", 3002, "Specify the port to listen to.")
	flag.StringVar(&elasticsearchUrl, "elasticsearchUrl", "http://127.0.0.1:9200", "Specify elasticsearch host and port.")
	flag.Parse()

	opts := getOpts()

	// Setup mux
	mux := http.NewServeMux()
	mux.Handle("/freelancer", opts.Handler(http.HandlerFunc(freelancer)))
	mux.Handle("/job", opts.Handler(http.HandlerFunc(job)))

	panic(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}

func job(w http.ResponseWriter, r *http.Request) {
	elasticSearchResponse, err := getElasticSearchResponse("job")
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}
	respond.With(w, r, http.StatusOK, buildSearchResponse(elasticSearchResponse, nil))
}

func freelancer(w http.ResponseWriter, r *http.Request) {
	elasticSearchResponse, err := getElasticSearchResponse("freelancer")
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}
	respond.With(w, r, http.StatusOK, buildSearchResponse(elasticSearchResponse, nil))
}

func getElasticSearchResponse(esType string) (ElasticSearchResponse, error) {
	var elasticSearchResponse ElasticSearchResponse
	resp, err := http.Get(elasticsearchUrl + "/fairlance/" + esType + "/_search")
	if err != nil {
		return elasticSearchResponse, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return elasticSearchResponse, err
	}

	err = json.Unmarshal(body, &elasticSearchResponse)
	if err != nil {
		return elasticSearchResponse, err
	}
	return elasticSearchResponse, nil
}

func buildSearchResponse(eResp ElasticSearchResponse, fn processFunc) interface{} {
	var entities = make([]interface{}, len(eResp.Hits.Hits))
	for key, hit := range eResp.Hits.Hits {
		if fn != nil {
			entities[key] = fn(hit.Source)
		} else {
			entities[key] = hit.Source
		}
	}

	response := struct {
		Total      int           `json:"total"`
		Collection []interface{} `json:"collection"`
	}{
		eResp.Hits.Total,
		entities,
	}

	return response
}

func getOpts() *respond.Options {
	return &respond.Options{
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
