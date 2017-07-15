package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	request, err := http.NewRequest("POST", "http://local.fairlance.io:8000/api/payment/deposit", bytes.NewBuffer([]byte(
		`{
			"projectID": 1
		}`,
	)))
	if err != nil {
		log.Fatal(err)
	}
	do(request)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	reader.ReadString('\n')
	request, err = http.NewRequest("POST", "http://local.fairlance.io:8000/api/payment/execute", bytes.NewBuffer([]byte(
		`{
			"projectID": 1
		}`,
	)))
	if err != nil {
		log.Fatal(err)
	}
	do(request)
}

func do(request *http.Request) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("status: %s, body: %s", response.Status, content)
		log.Fatal(err)
	}
	log.Printf("%s", content)
}
