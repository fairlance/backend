package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	reg "fairlance.io/registration"
)

type TestMailer struct{}

func (tm TestMailer) SendWelcomeMessage(email string) (string, error) {
	return "", nil
}

func buildTestContext(db string) *reg.RegistrationContext {
	// Setup context
	context := reg.NewContext(db)

	// override
	context.Mailer = TestMailer{}
	context.Logger = log.New(ioutil.Discard, "", 0)

	return context
}

func TestIndexHandler(t *testing.T) {
	setUp()
	req := getGETRequest()
	w := httptest.NewRecorder()
	reg.IndexHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusOK)
	if w.Body.String() != "[]\n" {
		t.Error(fmt.Printf("Body not ok (%q)\n", w.Body.String()))
	}
}

func TestPOSTIndexHandler(t *testing.T) {
	setUp()
	req := getPOSTRequest(nil)
	w := httptest.NewRecorder()
	reg.IndexHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusMethodNotAllowed)
	body := getBodyMap(w)
	if body["error"] != "Method not allowed! Use GET" && body["created"] != "" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func TestGETRegisterHandler(t *testing.T) {
	setUp()
	req := getGETRequest()
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusMethodNotAllowed)
	body := getBodyMap(w)
	if body["error"] != "Method not allowed! Use POST" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func TestRegisterHandlerForm(t *testing.T) {
	setUp()
	req := getPOSTRequest(nil)
	req.PostForm = url.Values{}
	req.PostForm.Set("email", "test@email.com")
	req.Header.Del("Content-Type")
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusCreated)
	body := getBodyMap(w)
	if _, present := body["created"]; !present || body["email"] != "test@email.com" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func TestRegisterHandlerJSON(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusCreated)
	body := getBodyMap(w)
	if _, present := body["created"]; !present || body["email"] != "test@email.com" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func TestAddingExistingUser(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	// body buffer gets emptied after RegisterHandler finishes
	// so we create a new request, with new body
	req = getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w = httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusConflict)
	body := getBodyMap(w)
	if body["error"] != "Email exists!" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func TestAddingEmptyUser(t *testing.T) {
	setUp()
	body := bytes.NewBuffer([]byte(`{"email":""}`))
	req := getPOSTRequest(body)
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusBadRequest)

	respBody := getBodyMap(w)
	if respBody["error"] != "Email missing!" {
		t.Error(fmt.Printf("Body not ok (%q)\n", respBody))
	}
}

func TestAddingInvalidJSON(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"invalid json`)))
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusBadRequest)
	body := getBodyMap(w)
	if body["error"] != "Request not valid JSON!" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func TestAddingInvalidUser(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"notemail.com"}`)))
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusBadRequest)
	body := getBodyMap(w)
	if body["error"] != "Email not valid!" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func TestAddingAndReadingRegisteredUser(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w := httptest.NewRecorder()
	reg.RegisterHandler(buildTestContext("registration_test"), w, req)

	req = getGETRequest()
	w = httptest.NewRecorder()
	reg.IndexHandler(buildTestContext("registration_test"), w, req)

	assertCode(t, w, http.StatusOK)
	var body []map[string]interface{}
	if err := json.Unmarshal([]byte(getBody(w)), &body); err != nil {
		t.Error(err)
	}
	if body[0]["email"] != "test@email.com" {
		t.Error(fmt.Printf("Body not ok (%q)\n", body))
	}
}

func setUp() {
	buildTestContext("registration_test").RegisteredUserRepository.RemoveAll()
}

// helper functions

func assertCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	if w.Code != expectedCode {
		t.Error(fmt.Printf("Code not ok (%d)\n", w.Code))
	}
}

func getBody(w *httptest.ResponseRecorder) string {
	return strings.Replace(w.Body.String(), "\n", "", -1)
}

func getBodyMap(w *httptest.ResponseRecorder) map[string]interface{} {
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(getBody(w)), &body); err != nil {
		panic(err)
	}

	return body
}

func getPOSTRequest(body io.Reader) *http.Request {
	req := getRequest("POST", body)
	req.Header.Set("Content-Type", "application/json")

	return req
}

func getGETRequest() *http.Request {
	return getRequest("GET", nil)
}

func getRequest(method string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, "http://example.com/foo", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Method = method

	return req
}
