package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	password   = flag.String("password", "", "Namecheap password.")
	dnsTimeout = flag.Int("timeout", 20, "Namecheap timeout.")
	client     = &http.Client{}
	ipURL      = "http://ipinfo.io/ip"
	updateURL  = "https://dynamicdns.park-your-domain.com/update?host=pi&domain=fairlance.io&password="
	fileName   = "lastDNSUpdate"
)

func get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	responseText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return responseText, nil
}

// <interface-response>
// 	<Command>SETDNSHOST</Command>
// 	<Language>eng</Language>
// 	<IP>127.0.0.1</IP>
// 	<ErrCount>0</ErrCount>
// 	<ResponseCount>0</ResponseCount>
// 	<Done>true</Done>
// 	<debug><![CDATA[]]></debug>
// </interface-response>

type NamecheapStatus struct {
	XMLName  xml.Name `xml:"interface-response"`
	ErrCount ErrCount
}

type ErrCount struct {
	XMLName xml.Name `xml:"ErrCount"`
	Count   string   `xml:",chardata"`
}

func updateIP() string {
	newIP, err := get(ipURL)
	if err != nil {
		return err.Error()
	}

	stringNewIP := strings.TrimSpace(string(newIP))

	updateStatus, err := get(updateURL + *password + "&ip=" + stringNewIP)
	if err != nil {
		return err.Error()
	}

	var namecheapStatus NamecheapStatus
	err = xml.Unmarshal(updateStatus, &namecheapStatus)
	if err != nil {
		return err.Error()
	}

	return "Update " + stringNewIP + ",  error count: " + namecheapStatus.ErrCount.Count
}

func logLastDNSUpdate(data string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	path := dir + "/" + fileName

	// detect if file exists
	_, err = os.Stat(path)

	// create file if not exists
	if !os.IsNotExist(err) {
		var err = os.Remove(path)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	ticker := time.NewTicker(time.Minute * time.Duration(*dnsTimeout))
	for t := range ticker.C {
		logLastDNSUpdate(fmt.Sprintf("Latest update", t, updateIP()))
	}
}
