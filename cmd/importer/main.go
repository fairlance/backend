package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/importer"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Indexed 50000 documents, in 6334.31s (average 126.69ms/doc)
func main() {
	var port = os.Getenv("PORT")
	var dbHost = os.Getenv("DB_HOST")
	var dbName = os.Getenv("DB_NAME")
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASS")
	var searcherURL = os.Getenv("SEARCHER_URL")
	var applicationURL = os.Getenv("APPLICATION_URL")

	// start the HTTP server
	http.Handle("/", importer.NewServeMux(importer.Options{
		DBHost:         dbHost,
		DBName:         dbName,
		DBUser:         dbUser,
		DBPass:         dbPass,
		SearcherURL:    searcherURL,
		ApplicationURL: applicationURL,
	}))

	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
