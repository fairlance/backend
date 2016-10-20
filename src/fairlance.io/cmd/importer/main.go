package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"fairlance.io/importer"

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
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// func batchIndex(i bleve.Index, docs map[string]interface{}) error {
// 	fmt.Println("Indexing...")
// 	var err error
// 	count := 0
// 	startTime := time.Now()
// 	batch := i.NewBatch()
// 	batchCount := 0
// 	for id, doc := range docs {
// 		batch.Index(id, doc)
// 		batchCount++

// 		if batchCount >= 100 {
// 			err = i.Batch(batch)
// 			if err != nil {
// 				return err
// 			}
// 			batch = i.NewBatch()
// 			batchCount = 0
// 		}
// 		count++
// 		if count%1000 == 0 {
// 			indexDuration := time.Since(startTime)
// 			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
// 			timePerDoc := float64(indexDuration) / float64(count)
// 			fmt.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
// 		}
// 	}
// 	// flush the last batch
// 	if batchCount > 0 {
// 		err = i.Batch(batch)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	indexDuration := time.Since(startTime)
// 	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
// 	timePerDoc := float64(indexDuration) / float64(count)
// 	fmt.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
// 	return nil
// }
