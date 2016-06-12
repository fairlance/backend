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
	Source Freelancer `json:"_source"`
}
type Freelancer struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func main() {
	flag.IntVar(&port, "port", 3002, "Specify the port to listen to.")
	flag.StringVar(&elasticsearchUrl, "elasticsearchUrl", "http://127.0.0.1:9200", "Specify elasticsearch host and port.")
	flag.Parse()

	opts := getOpts()

	// Setup mux
	mux := http.NewServeMux()
	mux.Handle("/", opts.Handler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp, err := http.Get(elasticsearchUrl + "/fairlance/freelancer/_search")
			if err != nil {
				respond.With(w, r, http.StatusBadGateway, err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				respond.With(w, r, http.StatusBadRequest, err)
				return
			}

			var elasticSearchResponse ElasticSearchResponse
			err = json.Unmarshal(body, &elasticSearchResponse)

			respond.With(w, r, resp.StatusCode, getSearchResponse(elasticSearchResponse))
		})),
	)

	panic(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}

func getSearchResponse(eResp ElasticSearchResponse) interface{} {
	var freelancers = make([]Freelancer, len(eResp.Hits.Hits))
	for key, hit := range eResp.Hits.Hits {
		freelancers[key] = hit.Source
	}

	response := struct {
		Total       int          `json:"total"`
		Freelancers []Freelancer `json:"freelancers"`
	}{
		eResp.Hits.Total,
		freelancers,
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
