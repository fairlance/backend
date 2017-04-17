package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Indexer interface {
	Index(index, docID string, document interface{}) error
	Delete(index, docID string) error
}

func NewIndexer(searcherURL string) *HTTPIndexer {
	return &HTTPIndexer{"http://" + searcherURL}
}

type HTTPIndexer struct {
	SearcherURL string
}

func (i *HTTPIndexer) Index(index, docID string, document interface{}) error {
	docBody, err := json.Marshal(document)
	if err != nil {
		return err
	}
	req, err := i.newRequest("PUT", index, docID, bytes.NewBuffer(docBody))
	if err != nil {
		return err
	}
	return i.doRequest(req)
}

func (i *HTTPIndexer) Delete(index, docID string) error {
	req, err := i.newRequest("DELETE", index, docID, nil)
	if err != nil {
		return err
	}
	return i.doRequest(req)
}

func (i *HTTPIndexer) newRequest(method, index, docID string, document io.Reader) (*http.Request, error) {
	url := i.SearcherURL + "/api/" + index + "/" + docID
	req, err := http.NewRequest(method, url, document)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (i *HTTPIndexer) doRequest(req *http.Request) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not execute the request, status: %s", resp.Status)
	}
	return nil
}
