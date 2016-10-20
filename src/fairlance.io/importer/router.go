package importer

import (
	"log"

	"github.com/gorilla/mux"
)

func NewRouter(options Options) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	db, err := getDB(options)
	if err != nil {
		log.Fatal(err)
	}
	router.Handle("/", indexHandler{
		options: options,
		db:      db,
	}).Methods("GET")

	return router
}
