package importer

import (
	"net/http"

	"encoding/json"

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
		doc, err := getDocFromDB(*i.db, pageState.Type, pageState.DB.DocID)
		if err != nil {
			pageState.Message = err.Error()
		}
		pageState.DB.Document = doc
	case "import":
		err := importDoc(*i.db, i.options, pageState.Type, pageState.DB.DocID)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "remove":
		err := deleteDocFromSearchEngine(i.options, pageState.Type, pageState.DB.DocID)
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
		entities, err := doSearch(i.options, pageState)
		if err != nil {
			pageState.Message = err.Error()
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
		pageState.Search.Tags, err = getSearchTags(i.options)
		if err != nil {
			pageState.Message = err.Error()
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pageState)
}
