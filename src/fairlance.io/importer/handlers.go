package importer

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

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
		err := doIndex(i.options, *i.db, pageState.Type)
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

func newPage(r *http.Request) page {
	pageState := page{Message: "ok"}
	query := r.URL.Query()
	offset := 0
	if query.Get("offset") != "" {
		o, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {
			pageState.Message = err.Error()
		}
		offset = int(o)
	}
	pageState.Offset = offset

	limit := 10
	if query.Get("limit") != "" {
		l, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if err != nil {
			pageState.Message = err.Error()
		}
		limit = int(l)
	}
	pageState.Limit = limit

	pageState.Type = "jobs"
	if query.Get("type") != "" {
		pageState.Type = query.Get("type")
	}

	return pageState
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
