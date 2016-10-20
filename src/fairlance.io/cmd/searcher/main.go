//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/blevesearch/bleve"
	bleveHttp "github.com/blevesearch/bleve/http"

	// import general purpose configuration
	_ "github.com/blevesearch/bleve/config"
)

var port = flag.String("port", "3003", "http listen address")
var dataDir = flag.String("dataDir", "/tmp/indices", "data directory")

func init() {
	f, err := os.OpenFile("/var/log/fairlance/searcher.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	flag.Parse()

	// walk the data dir and register index names
	dirEntries, err := ioutil.ReadDir(*dataDir)
	if err != nil {
		log.Printf("Creating folder %s ...\n", *dataDir)
		os.MkdirAll(*dataDir, 0777)
		dirEntries, err = ioutil.ReadDir(*dataDir)
		if err != nil {
			log.Fatalf("error reading data dir: %v", err)
		}
	}

	if len(dirEntries) == 0 {
		if err := buildIndices(); err != nil {
			log.Fatalf("error reading data dir: %v", err)
		} else {
			log.Println("Created indices.")
		}
	}

	for _, dirInfo := range dirEntries {
		indexPath := *dataDir + string(os.PathSeparator) + dirInfo.Name()

		// skip single files in data dir since a valid index is a directory that
		// contains multiple files
		if !dirInfo.IsDir() {
			log.Printf("not registering %s, skipping", indexPath)
			continue
		}

		i, err := bleve.Open(indexPath)
		if err != nil {
			log.Printf("error opening index %s: %v", indexPath, err)
		} else {
			log.Printf("registered index: %s", dirInfo.Name())
			bleveHttp.RegisterIndexName(dirInfo.Name(), i)
			// set correct name in stats
			i.SetName(dirInfo.Name())
		}
	}

	router := mux.NewRouter()
	router.StrictSlash(true)

	createIndexHandler := bleveHttp.NewCreateIndexHandler(*dataDir)
	createIndexHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}", createIndexHandler).Methods("PUT")

	getIndexHandler := bleveHttp.NewGetIndexHandler()
	getIndexHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}", getIndexHandler).Methods("GET")

	deleteIndexHandler := bleveHttp.NewDeleteIndexHandler(*dataDir)
	deleteIndexHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}", deleteIndexHandler).Methods("DELETE")

	listIndexesHandler := bleveHttp.NewListIndexesHandler()
	router.Handle("/api", listIndexesHandler).Methods("GET")

	docIndexHandler := bleveHttp.NewDocIndexHandler("")
	docIndexHandler.IndexNameLookup = indexNameLookup
	docIndexHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}", docIndexHandler).Methods("PUT")

	docCountHandler := bleveHttp.NewDocCountHandler("")
	docCountHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_count", docCountHandler).Methods("GET")

	docGetHandler := bleveHttp.NewDocGetHandler("")
	docGetHandler.IndexNameLookup = indexNameLookup
	docGetHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}", docGetHandler).Methods("GET")

	docDeleteHandler := bleveHttp.NewDocDeleteHandler("")
	docDeleteHandler.IndexNameLookup = indexNameLookup
	docDeleteHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}", docDeleteHandler).Methods("DELETE")

	searchHandler := bleveHttp.NewSearchHandler("")
	searchHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_search", searchHandler).Methods("POST")

	listFieldsHandler := bleveHttp.NewListFieldsHandler("")
	listFieldsHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_fields", listFieldsHandler).Methods("GET")

	debugHandler := bleveHttp.NewDebugDocumentHandler("")
	debugHandler.IndexNameLookup = indexNameLookup
	debugHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}/_debug", debugHandler).Methods("GET")

	aliasHandler := bleveHttp.NewAliasHandler()
	router.Handle("/api/_aliases", aliasHandler).Methods("POST")

	// start the HTTP server
	log.Printf("listening on :%s\n", *port)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func buildIndices() error {
	for _, indexName := range []string{"jobs", "freelancers"} {
		log.Printf("creating %s\n", indexName)
		mapping := bleve.NewIndexMapping()
		_, err := bleve.New(*dataDir+"/"+indexName, mapping)
		if err != nil {
			return err
		}
	}

	return nil
}

func indexNameLookup(req *http.Request) string {
	return muxVariableLookup(req, "indexName")
}

func muxVariableLookup(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}

func docIDLookup(req *http.Request) string {
	return muxVariableLookup(req, "docID")
}
