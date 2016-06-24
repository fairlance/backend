package main_test

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

	app "fairlance.io/application"
)

var (
	appContext   *app.ApplicationContext
	emptyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

func TestIdHandler(t *testing.T) {
	is := is.New(t)

	r, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle("/{id}", app.IdHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "id").(uint)
		is.Equal(id, 1)
	}))).Methods("GET")
	router.ServeHTTP(w, r)
}

func TestUserHandler(t *testing.T) {
	is := is.New(t)
	requestBody := `
	{
	  "password": "123",
	  "email": "pera@gmail.com",
	  "firstName":"Pera",
	  "lastName":"Peric"
	}`

	w := httptest.NewRecorder()
	r := getRequest("GET", requestBody)
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	app.RegisterUserHandler(emptyHandler).ServeHTTP(w, r)
	user := context.Get(r, "user").(*app.User)
	is.Equal(user.FirstName, "Pera")
	is.Equal(user.LastName, "Peric")
	is.Equal(user.Email, "pera@gmail.com")
}

func TestUserHandlerWithInvalidBody(t *testing.T) {
	is := is.New(t)
	requestBody := `
	{
		"empty": "invalid body"
	}`

	w := httptest.NewRecorder()
	r := getRequest("GET", requestBody)
	app.RegisterUserHandler(emptyHandler).ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusBadRequest)
	var errorBody map[string]string
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &errorBody))
	is.OK(errorBody["Email"])
	is.OK(errorBody["FirstName"])
	is.OK(errorBody["LastName"])
	is.OK(errorBody["Password"])
}

func TestUserHandlerWithInvalidEmail(t *testing.T) {
	is := is.New(t)
	requestBody := `
	{
	  "email": "invalid email",
	  "password": "123",
	  "firstName":"Pera",
	  "lastName":"Peric"
	}`

	w := httptest.NewRecorder()
	r := getRequest("GET", requestBody)
	app.RegisterUserHandler(emptyHandler).ServeHTTP(w, r)
	is.Equal(w.Code, http.StatusBadRequest)
	var body map[string]string
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.OK(body["Email"])
}

func buildTestContext() *app.ApplicationContext {
	options := app.ContextOptions{"application_test", "fairlance", "fairlance", "secret"}
	context, err := app.NewContext(options)
	if err != nil {
		panic(err)
	}

	return context
}

func setUp() {
	appContext = buildTestContext()
	appContext.DropTables()
	appContext.CreateTables()
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
