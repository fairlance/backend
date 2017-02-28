package application

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

type testNotifier struct{}

func (n *testNotifier) Notify(not *notification) error { return nil }

func getRequest(appContext *ApplicationContext, requestBody string) *http.Request {
	req, err := http.NewRequest("GET", "http://fairlance.io/", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	appContext.Notifier = &testNotifier{}
	req.Header.Set("Content-Type", "application/json")
	context.Set(req, "context", appContext)

	return req
}
