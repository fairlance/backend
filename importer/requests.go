package importer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fairlance/backend/application"
	"github.com/jinzhu/gorm"
)

func deleteDocFromSearchEngine(options Options, index, docID string) error {
	url := options.SearcherURL + "/api/" + index + "/" + docID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

func getTotalInSearchEngine(options Options, index string) (int, error) {
	var response struct {
		Status string
		Count  int
	}
	url := options.SearcherURL + "/api/" + index + "/_count"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response.Count, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response.Count, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return response.Count, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response.Count, err
	}

	return response.Count, nil
}

func deleteAllFromSearchEngine(options Options, index string) error {
	url := options.SearcherURL + "/api/" + index

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	req, err = http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

func importDocs(options Options, index string, docs map[string]interface{}) error {
	count := 0
	startTime := time.Now()
	for id, doc := range docs {
		if err := importDocument(options, index, id, doc); err != nil {
			return err
		}
		count++
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}

func importDoc(db gorm.DB, options Options, index, docID string) error {
	id, err := strconv.ParseInt(docID, 10, 64)
	if err != nil {
		return err
	}

	var doc interface{}
	switch index {
	case "jobs":
		job := application.Job{}
		if err := db.Preload("Client").Find(&job, id).Error; err != nil {
			return err
		}
		doc = job
	case "freelancers":
		freelancer := application.Freelancer{}
		if err := db.Find(&freelancer, id).Error; err != nil {
			return err
		}
		doc = freelancer
	}
	if err := importDocument(options, index, docID, doc); err != nil {
		return err
	}
	return nil
}

func importDocument(options Options, index, docID string, doc interface{}) error {
	docBody, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	url := options.SearcherURL + "/api/" + index + "/" + docID
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(docBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// func batchIndex(i bleve.Index, docs map[string]interface{}) error {
// 	fmt.Println("Indexing...")
// 	var err error
// 	count := 0
// 	startTime := time.Now()
// 	batch := i.NewBatch()
// 	batchCount := 0
// 	for id, doc := range docs {
// 		batch.Index(id, doc)
// 		batchCount++

// 		if batchCount >= 100 {
// 			err = i.Batch(batch)
// 			if err != nil {
// 				return err
// 			}
// 			batch = i.NewBatch()
// 			batchCount = 0
// 		}
// 		count++
// 		if count%1000 == 0 {
// 			indexDuration := time.Since(startTime)
// 			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
// 			timePerDoc := float64(indexDuration) / float64(count)
// 			fmt.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
// 		}
// 	}
// 	// flush the last batch
// 	if batchCount > 0 {
// 		err = i.Batch(batch)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	indexDuration := time.Since(startTime)
// 	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
// 	timePerDoc := float64(indexDuration) / float64(count)
// 	fmt.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
// 	return nil
// }
