package importer

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(options Options) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	db, err := getDB(options)
	if err != nil {
		log.Fatal(err)
	}

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index").Parse(htmlTemplate)
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	})).Methods("GET")

	router.Handle("/json", indexHandlerJSON{
		options: options,
		db:      db,
	}).Methods("GET")
	router.Handle("/json", searchHandler{
		options: options,
	}).Methods("POST", "OPTIONS")

	router.Handle("/websockettest", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main := MustAsset("templates/websockettest.html")

		tmpl, err := template.New("messages").Parse(string(main))
		if err != nil {
			log.Fatal(err)
		}
		if err := tmpl.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
		return
	}))

	return router
}
