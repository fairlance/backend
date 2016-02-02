package registration

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type TestMailer struct{}

func (tm TestMailer) SendWelcomeMessage(email string) (string, error) {
	return "", nil
}

func buildTestContext(db string) *RegistrationContext {
	// Setup context
	context := NewContext("test")

	// override
	context.mailer = TestMailer{}
	context.Logger = log.New(ioutil.Discard, "", 0)

	return context
}

func TestIndexHandler(t *testing.T) {
	setUp()
	req := getGETRequest()
	w := httptest.NewRecorder()
	IndexHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusOK)
	assertBody(t, w, "[]")
}

func TestPOSTIndexHandler(t *testing.T) {
	setUp()
	req := getPOSTRequest(nil)
	w := httptest.NewRecorder()
	IndexHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusMethodNotAllowed)
	assertBody(t, w, "{\"error\":\"Method not allowed! Use GET\"}")
}

func TestGETRegisterHandler(t *testing.T) {
	setUp()
	req := getGETRequest()
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusMethodNotAllowed)
	assertBody(t, w, "{\"error\":\"Method not allowed! Use POST\"}")
}

func TestRegisterHandlerForm(t *testing.T) {
	setUp()
	req := getPOSTRequest(nil)
	req.PostForm = url.Values{}
	req.PostForm.Set("email", "test@email.com")
	req.Header.Del("Content-Type")
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusCreated)
	assertBody(t, w, "{\"email\":\"test@email.com\"}")
}

func TestRegisterHandlerJSON(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusCreated)
	assertBody(t, w, "{\"email\":\"test@email.com\"}")
}

func TestAddingExistingUser(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	// body buffer gets emptied after RegisterHandler finishes
	// so we create a new request, with new body
	req = getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w = httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusConflict)
	assertBody(t, w, "{\"error\":\"Email exists!\"}")
}

func TestAddingEmptyUser(t *testing.T) {
	setUp()
	body := bytes.NewBuffer([]byte(`{"email":""}`))
	req := getPOSTRequest(body)
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusBadRequest)
	assertBody(t, w, "{\"error\":\"Email missing!\"}")
}

func TestAddingInvalidJSON(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"invalid json`)))
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusBadRequest)
	assertBody(t, w, "{\"error\":\"Request not valid JSON!\"}")
}

func TestAddingInvalidUser(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"notemail.com"}`)))
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusBadRequest)
	assertBody(t, w, "{\"error\":\"Email not valid!\"}")
}

func TestAddingAndReadingRegisteredUser(t *testing.T) {
	setUp()
	req := getPOSTRequest(bytes.NewBuffer([]byte(`{"email":"test@email.com"}`)))
	w := httptest.NewRecorder()
	RegisterHandler(buildTestContext("test"), w, req)

	req = getGETRequest()
	w = httptest.NewRecorder()
	IndexHandler(buildTestContext("test"), w, req)

	assertCode(t, w, http.StatusOK)
	assertBody(t, w, "[{\"email\":\"test@email.com\"}]")
}

func setUp() {
	buildTestContext("test").userRepository.getUsers().RemoveAll(nil)
}

// helper functions

func assertCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	if w.Code != expectedCode {
		t.Error(fmt.Printf("Code not ok (%d)\n", w.Code))
	}
}

func assertBody(t *testing.T, w *httptest.ResponseRecorder, expectedBody string) {
	if strings.Replace(w.Body.String(), "\n", "", -1) != expectedBody {
		t.Error(fmt.Printf("Body not ok (%q)\n", w.Body.String()))
	}
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
