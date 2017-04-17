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

type testIndexerIndexCallback func(index, docID string, document interface{}) error
type testIndexerDeleteCallback func(index, docID string) error

type testIndexer struct {
	indexCallback  testIndexerIndexCallback
	deleteCallback testIndexerDeleteCallback
}

func (i *testIndexer) Index(index, docID string, document interface{}) error {
	return i.indexCallback(index, docID, document)
}

func (i *testIndexer) Delete(index, docID string) error {
	return i.deleteCallback(index, docID)
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
	if appContext.Indexer == nil {
		appContext.Indexer = &testIndexer{
			indexCallback:  func(index, docID string, document interface{}) error { return nil },
			deleteCallback: func(index, docID string) error { return nil },
		}
	}
	req.Header.Set("Content-Type", "application/json")
	context.Set(req, "context", appContext)

	return req
}
