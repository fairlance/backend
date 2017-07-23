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
	log.SetFlags(log.Lshortfile)
	var port = os.Getenv("PORT")
	// start the HTTP server
	http.Handle("/", importer.NewServeMux(importer.Options{
		DBHost:         os.Getenv("DB_HOST"),
		DBName:         os.Getenv("DB_NAME"),
		DBUser:         os.Getenv("DB_USER"),
		DBPass:         os.Getenv("DB_PASS"),
		SearcherURL:    os.Getenv("SEARCHER_URL"),
		ApplicationURL: os.Getenv("APPLICATION_URL"),
	}))
	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
