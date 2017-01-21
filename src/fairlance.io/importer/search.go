package importer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// var baseURL = "http://localhost:3002"

func getDocFromSearchEngine(options Options, index, docID string) (map[string]interface{}, error) {
	url := options.SearcherURL + "/api/" + index + "/" + docID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var doc map[string]interface{}
	err = json.Unmarshal(body, &doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
