package importer

import (
	"html/template"
	"log"
	"net/http"
)

func NewServeMux(options Options) *http.ServeMux {
	mux := http.NewServeMux()
	db, err := getDB(options)
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main := MustAsset("templates/index.html")
		tmpl, err := template.New("index").Parse(string(main))
		if err != nil {
			log.Fatal(err)
		}
		if err := tmpl.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
		return
	}))
	mux.Handle("/json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			indexHandlerJSON{options: options, db: db}.ServeHTTP(w, r)
			return
		}
		searchHandler{options: options}.ServeHTTP(w, r)
	}))
	mux.Handle("/websockettest", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main := MustAsset("templates/websockettest.html")
		tmpl, err := template.New("websockettest").Parse(string(main))
		if err != nil {
			log.Fatal(err)
		}
		if err := tmpl.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
		return
	}))
	mux.Handle("/payment", paymentHandler(db.DB()))
	return mux
}
