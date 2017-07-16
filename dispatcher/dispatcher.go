package dispatcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func doGET(client *http.Client, url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("status: %s, body: %s, url: %s", response.Status, content, url)
		log.Printf("could not execute get request: %v", err)
		return content, err
	}
	return content, nil
}

func doPOST(client *http.Client, url string, b []byte) error {
	request, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		err = fmt.Errorf("status: %s, body: %s, url: %s", response.Status, content, url)
		log.Printf("could not execute post request: %v", err)
		return err
	}
	return nil
}
