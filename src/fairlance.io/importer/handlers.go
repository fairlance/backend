package importer

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

type indexHandler struct {
	options Options
	db      *gorm.DB
}

func (i indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pageState := newPage(r)
	query := r.URL.Query()
	action := query.Get("action")
	switch action {
	case "import_all":
		err := doImport(i.options, *i.db, pageState.Type)
		if err != nil {
			pageState.Message = err.Error()
		}
	case "get":
		doc, err := getDocFromSearchEngine(i.options, pageState.Type, query.Get("docID"))
		if err != nil {
			pageState.Message = err.Error()
		}
		pageState.Document = doc
	case "import":
		err := importDoc(*i.db, i.options, pageState.Type, query.Get("docID"))
		if err != nil {
			pageState.Message = err.Error()
		}
	case "remove":
		err := deleteDocFromSearchEngine(i.options, pageState.Type, query.Get("docID"))
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
	}

	var err error
	switch pageState.Type {
	case "jobs":
		pageState.Entities, pageState.TotalInDB, err = getJobsFromDB(*i.db, pageState.Offset, pageState.Limit)
	case "freelancers":
		pageState.Entities, pageState.TotalInDB, err = getFreelancersFromDB(*i.db, pageState.Offset, pageState.Limit)
	}
	if err != nil {
		pageState.Message = err.Error()
	}

	pageState.TotalInSearchEngine, err = getTotalInSearchEngine(i.options, pageState.Type)
	if err != nil {
		pageState.Message = err.Error()
	}

	renderTemplate(w, pageState)
}

func renderTemplate(w http.ResponseWriter, data page) {
	t, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
