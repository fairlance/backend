package main

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/cheekybits/is"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"net/http"
	"net/http/httptest"
)

var (
	appContext   *ApplicationContext
	emptyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

func TestIndex(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	Index(w, r)

	is.Equal(w.Code, http.StatusOK)
	var data string
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, "Hi")
}

func TestIdHandler(t *testing.T) {
	is := is.New(t)

	r, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle("/{id}", IdHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "id").(uint64)
		is.Equal(id, 1)
	}))).Methods("GET")
	router.ServeHTTP(w, r)
}

func buildTestContext(db string) *ApplicationContext {
	context, err := NewContext(db)
	if err != nil {
		panic(err)
	}

	return context
}

func setUp() {
	appContext = buildTestContext("application_test")
	appContext.TruncateTables()
}

func getRequest(method string, requestBody string) *http.Request {
	req, err := http.NewRequest(method, "http://fairlance.io/", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	context.Set(req, "context", appContext)

	return req
}
