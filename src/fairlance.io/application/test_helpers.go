package application

import (
	"bytes"
	"log"
	"net/http"

	"fairlance.io/notifier"

	"github.com/gorilla/context"
)

type testNotifierCallback func(notification *notifier.Notification) error

type testNotifier struct {
	callback testNotifierCallback
}

func (n *testNotifier) Notify(notification *notifier.Notification) error {
	return n.callback(notification)
}

func getRequest(appContext *ApplicationContext, requestBody string) *http.Request {
	req, err := http.NewRequest("GET", "http://fairlance.io/", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	if appContext.Notifier == nil {
		appContext.Notifier = &testNotifier{
			callback: func(notification *notifier.Notification) error { return nil },
		}
	}
	req.Header.Set("Content-Type", "application/json")
	context.Set(req, "context", appContext)

	return req
}
