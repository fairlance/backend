package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var baseURL = "http://localhost:3002"

func getSearchTags(options Options) ([]string, error) {
	url := baseURL + "/job/tags"

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

	type resultStruct struct {
		Data struct {
			Tags []string `json:"tags"`
		} `json:"data"`
	}
	var result = resultStruct{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Tags, nil
}

func doSearch(options Options, pageState page) (map[string]interface{}, error) {
	period := ""
	if pageState.Search.Period != "" {
		period = "period=" + pageState.Search.Period
	}
	priceFrom := ""
	if pageState.Search.PriceFrom != "" {
		priceFrom = "&price_from=" + pageState.Search.PriceFrom
	}
	priceTo := ""
	if pageState.Search.PriceTo != "" {
		priceTo = "&price_to=" + pageState.Search.PriceTo
	}
	tags := ""
	for _, tag := range pageState.Search.Tags {
		tags = tags + "&tags=" + tag
	}
	index := "job"
	url := fmt.Sprintf("%s/%s?%s%s%s%s", baseURL, index, period, priceFrom, priceTo, tags)

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

	type resultStruct struct {
		Data struct {
			Items []map[string]interface{} `json:"items"`
		} `json:"data"`
	}
	var result = resultStruct{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var items = make(map[string]interface{})
	for _, item := range result.Data.Items {
		id := strconv.FormatUint(uint64(item["id"].(float64)), 10)
		items[id] = item
	}

	return items, nil
}
