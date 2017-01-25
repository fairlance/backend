package importer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"strings"

	"github.com/jinzhu/gorm"
)

type indexHandlerJSON struct {
	options Options
	db      *gorm.DB
}

func (i indexHandlerJSON) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pageState := newPage(r)
	switch pageState.Action {
	case "import_all":
		err := doImport(i.options, *i.db, pageState.Type)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "get":
		doc, err := getDocFromDB(*i.db, pageState.Type, pageState.DocID)
		if err != nil {
			pageState.Message = err.Error()
		}
		pageState.DB.Document = doc
	case "import":
		err := importDoc(*i.db, i.options, pageState.Type, pageState.DocID)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "remove":
		err := deleteDocFromSearchEngine(i.options, pageState.Type, pageState.DocID)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "re_generate_test_data":
		err := reGenerateTestData(*i.db, pageState.Type)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "delete_all_from_db":
		err := deleteAllFromDB(*i.db, pageState.Type)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "delete_all_from_search_engine":
		err := deleteAllFromSearchEngine(i.options, pageState.Type)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "search":
		entity, err := getDocFromSearchEngine(i.options, pageState.Type, pageState.DocID)
		if err != nil {
			pageState.Message = err.Error()
		}
		entities := make(map[string]interface{})
		id, ok := entity["id"].(string)
		if ok {
			entities[id] = entity
		}
		pageState.Entities = entities
	}

	var err error
	switch pageState.Tab {
	case "db":
		switch pageState.Type {
		case "jobs":
			pageState.Entities, pageState.DB.TotalInDB, err = getJobsFromDB(*i.db, pageState.Offset, pageState.Limit)
		case "freelancers":
			pageState.Entities, pageState.DB.TotalInDB, err = getFreelancersFromDB(*i.db, pageState.Offset, pageState.Limit)
		}
		pageState.DB.TotalInSearchEngine, err = getTotalInSearchEngine(i.options, pageState.Type)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "search":

	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pageState)
}

type searchHandler struct {
	options Options
}

func (handler searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type responseStruct struct {
		RawData map[string]interface{}
		Message string
	}

	response := responseStruct{Message: "ok"}

	decoder := json.NewDecoder(r.Body)
	var body map[string]interface{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Println(err)
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := handler.options.SearcherURL + "/api/" + body["url"].(string)
	log.Println(url)
	req, err := http.NewRequest(strings.ToUpper(body["method"].(string)), url, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Println(url)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = fmt.Sprintf("Status %d", resp.StatusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	var doc map[string]interface{}
	err = json.Unmarshal(b, &doc)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	response.RawData = doc

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
