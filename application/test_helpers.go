package application

import (
	"bytes"
	"log"
	"net/http"

	"github.com/fairlance/backend/dispatcher"

	"github.com/gorilla/context"
)

type testNotifierCallback func(notification *dispatcher.Notification) error

type testNotifier struct {
	callback testNotifierCallback
}

func (n *testNotifier) Notify(notification *dispatcher.Notification) error {
	return n.callback(notification)
}

type testMessagingCallback func(meassage *dispatcher.Message) error

type testMessaging struct {
	callback testMessagingCallback
}

func (n *testMessaging) Send(message *dispatcher.Message) error {
	return n.callback(message)
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

type testPaymentCallback func(projectID uint) error

type testPayment struct {
	callback testPaymentCallback
}

func (p *testPayment) Deposit(projectID uint) error {
	return p.callback(projectID)
}

func (p *testPayment) Execute(projectID uint) error {
	return p.callback(projectID)
}

func getRequest(appContext *ApplicationContext, requestBody string) *http.Request {
	req, err := http.NewRequest("GET", "http://github.com/fairlance/", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatal(err)
	}
	if appContext.NotificationDispatcher == nil {
		appContext.NotificationDispatcher = NewNotificationDispatcher(&testNotifier{
			callback: func(notification *dispatcher.Notification) error { return nil },
		})
	}
	if appContext.MessagingDispatcher == nil {
		appContext.MessagingDispatcher = NewMessagingDispatcher(&testMessaging{
			callback: func(message *dispatcher.Message) error { return nil },
		})
	}
	if appContext.PaymentDispatcher == nil {
		appContext.PaymentDispatcher = NewPaymentDispatcher(&testPayment{
			callback: func(projectID uint) error { return nil },
		})
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
