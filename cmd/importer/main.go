package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/fairlance/backend/importer"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	port        int
	dbHost      string
	dbName      string
	dbUser      string
	dbPass      string
	searcherURL string
)

func init() {
	// f, err := os.OpenFile("/var/log/fairlance/importer.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// log.SetOutput(f)
}

// Indexed 50000 documents, in 6334.31s (average 126.69ms/doc)
func main() {
	flag.IntVar(&port, "port", 3004, "http listen address")
	flag.StringVar(&dbHost, "dbHost", "localhost", "DB host.")
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "fairlance", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "fairlance", "Db user's password.")
	flag.StringVar(&searcherURL, "searcherURL", "http://localhost:3003", "Url of the searcher.")
	flag.Parse()

	// start the HTTP server
	http.Handle("/", importer.NewRouter(importer.Options{
		DBHost:      dbHost,
		DBName:      dbName,
		DBUser:      dbUser,
		DBPass:      dbPass,
		SearcherURL: searcherURL,
	}))

	log.Printf("Listening on: %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
