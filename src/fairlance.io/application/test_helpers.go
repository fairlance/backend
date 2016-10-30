package application

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

func getRequest(appContext *ApplicationContext, requestBody string) *http.Request {
	req, err := http.NewRequest("GET", "http://fairlance.io/", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	context.Set(req, "context", appContext)

	return req
}
