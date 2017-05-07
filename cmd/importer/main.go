package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/importer"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	dbName      string
	dbUser      string
	dbPass      string
	searcherURL string
)

var port = flag.String("port", "3004", "http listen address")

func init() {
	f, err := os.OpenFile("/var/log/fairlance/importer.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

// Indexed 50000 documents, in 6334.31s (average 126.69ms/doc)
func main() {
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "fairlance", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "fairlance", "Db user's password.")
	flag.StringVar(&searcherURL, "searcherURL", "http://localhost:3003", "Url of the searcher.")
	flag.Parse()

	// start the HTTP server
	http.Handle("/", importer.NewRouter(importer.Options{
		DBName:      dbName,
		DBUser:      dbUser,
		DBPass:      dbPass,
		SearcherURL: searcherURL,
	}))
	http.Handle("/websockettest", importer.NewMessagesHandler())
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
