package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	setUp()
	req := getRequest("GET")
	w := httptest.NewRecorder()
	IndexHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusOK)
	assertBody(t, w, "[]")
}

func TestPOSTIndexHandler(t *testing.T) {
	setUp()
	req := getRequest("POST")
	w := httptest.NewRecorder()
	IndexHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusForbidden)
	assertBody(t, w, "{\"error\":\"Method not allowed! Use GET\"}")
}

func TestRegisterHandler(t *testing.T) {
	setUp()
	req := getRequest("POST")
	req.PostForm.Set("email", "test@email.com")
	req.ParseForm()
	w := httptest.NewRecorder()
	RegisterHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusCreated)
	assertBody(t, w, "{\"email\":\"test@email.com\"}")
}

func TestGETRegisterHandler(t *testing.T) {
	setUp()
	req := getRequest("GET")
	w := httptest.NewRecorder()
	RegisterHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusForbidden)
	assertBody(t, w, "{\"error\":\"Method not allowed! Use POST\"}")
}

func TestAddingExistingUser(t *testing.T) {
	setUp()
	req := getRequest("POST")
	req.PostForm.Set("email", "test@email.com")
	req.ParseForm()
	w := httptest.NewRecorder()
	RegisterHandler(buildContext("test"), w, req)

	w = httptest.NewRecorder()
	RegisterHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusConflict)
	assertBody(t, w, "{\"error\":\"Email exists!\"}")
}

func TestAddingEmptyUser(t *testing.T) {
	setUp()
	req := getRequest("POST")
	req.PostForm.Set("email", "")
	req.ParseForm()
	w := httptest.NewRecorder()
	RegisterHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusBadRequest)
	assertBody(t, w, "{\"error\":\"Email missing!\"}")
}

func TestAddingInvalidUser(t *testing.T) {
	setUp()
	req := getRequest("POST")
	req.PostForm.Set("email", "notanemail.com")
	req.ParseForm()
	w := httptest.NewRecorder()
	RegisterHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusBadRequest)
	assertBody(t, w, "{\"error\":\"Email not valid!\"}")
}

func TestAddingAndReadingRegisteredUser(t *testing.T) {
	setUp()
	req := getRequest("POST")
	req.PostForm.Set("email", "test@email.com")
	req.ParseForm()
	w := httptest.NewRecorder()
	RegisterHandler(buildContext("test"), w, req)

	req = getRequest("GET")
	w = httptest.NewRecorder()
	IndexHandler(buildContext("test"), w, req)

	assertCode(t, w, http.StatusOK)
	assertBody(t, w, "[{\"email\":\"test@email.com\"}]")
}

func assertCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	if w.Code != expectedCode {
		t.Error(fmt.Printf("Code not ok (%s)", w.Code))
	}
}

func assertBody(t *testing.T, w *httptest.ResponseRecorder, expectedBody string) {
	if strings.Replace(w.Body.String(), "\n", "", -1) != expectedBody {
		t.Error(fmt.Printf("Body not ok (%q)", w.Body.String()))
	}
}

func setUp() {
	buildContext("test").userRepository.getUsers().RemoveAll(nil)
}

func getRequest(method string) *http.Request {
	req, err := http.NewRequest(method, "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.PostForm = url.Values{}
	req.Method = method

	return req
}
