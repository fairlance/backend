package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"fairlance.io/application"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
)

var (
	dbName     string
	dbUser     string
	dbPass     string
	indicesDir string
)

// Indexed 50000 documents, in 6334.31s (average 126.69ms/doc)
func main() {
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "fairlance", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "fairlance", "Db user's password.")
	flag.StringVar(&indicesDir, "indicesDir", "/tmp/indices", "Location where the indices are located.")
	flag.Parse()

	jobsIndex, err := getIndex("jobs")
	if err != nil {
		log.Fatal(err)
	}

	freelancersIndex, err := getIndex("freelancers")
	if err != nil {
		log.Fatal(err)
	}

	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}

	jobsFromDB, err := jobs(db)
	if err != nil {
		log.Fatal(err)
	}

	if err = batchIndex(jobsIndex, jobsFromDB); err != nil {
		log.Fatal(err)
	}

	freelancersFromDB, err := freelancers(db)
	if err != nil {
		log.Fatal(err)
	}

	if err = batchIndex(freelancersIndex, freelancersFromDB); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}

func jobs(db *gorm.DB) (map[string]interface{}, error) {

	db.DropTableIfExists(&application.Job{}, &application.Tag{})
	db.CreateTable(&application.Job{}, &application.Tag{})

	// PostgreSQL only supports 65535 parameters
	for i := 0; i < 100; i++ {
		db.Create(&application.Job{
			Name:      fmt.Sprintf("Job %d", i),
			Summary:   fmt.Sprintf("Job Summary %d", i),
			Details:   fmt.Sprintf("Job Description %d", i),
			ClientID:  1,
			Price:     123*i%200 + 200,
			StartDate: time.Now().Add(time.Duration(i*24+1) * time.Hour),
			Tags: []application.Tag{
				application.Tag{Tag: fmt.Sprintf("tag_%d", i)},
				application.Tag{Tag: fmt.Sprintf("tag_%d", i+i)},
			},
		})
	}
	fmt.Println("Done Faking.")

	jobs := []application.Job{}
	jobsMap := make(map[string]interface{})
	if err := db.Preload("Tags").Preload("Client").Find(&jobs).Error; err != nil {
		return jobsMap, err
	}

	for _, job := range jobs {
		id := strconv.FormatUint(uint64(job.ID), 10)
		jobsMap[id] = job
	}

	fmt.Printf("Found %d jobs ...\n", len(jobs))
	return jobsMap, nil
}

func freelancers(db *gorm.DB) (map[string]interface{}, error) {
	freelancersMap := make(map[string]interface{})
	freelancers := []application.Freelancer{}
	if err := db.Preload("Skills", "owner_type = ?", "freelancers").Find(&freelancers).Error; err != nil {
		return freelancersMap, err
	}

	for _, freelancer := range freelancers {
		id := strconv.FormatUint(uint64(freelancer.ID), 10)
		freelancersMap[id] = freelancer
	}

	fmt.Printf("Found %d freelancers ...\n", len(freelancers))
	return freelancersMap, nil
}

func getIndex(dbName string) (bleve.Index, error) {
	fmt.Printf("Opening %s index ...\n", dbName)
	index, err := bleve.Open(indicesDir + "/" + dbName)
	if err == bleve.ErrorIndexPathDoesNotExist {
		fmt.Printf("%s index not found. Creating ...\n", dbName)
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indicesDir+"/"+dbName, mapping)
		if err != nil {
			return index, err
		}
	} else if err != nil {
		return index, err
	} else {
		fmt.Printf("Opening existing index...")
	}

	fmt.Printf("Opened %s index\n", dbName)
	return index, nil
}

func getDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", dbName, dbUser, dbPass))
	if err != nil {
		fmt.Println(err.Error())
		return db, err
	}

	fmt.Println("Opened DB")
	return db, nil
}

func batchIndex(i bleve.Index, docs map[string]interface{}) error {
	fmt.Println("Indexing...")
	var err error
	count := 0
	startTime := time.Now()
	batch := i.NewBatch()
	batchCount := 0
	for id, doc := range docs {
		batch.Index(id, doc)
		batchCount++

		if batchCount >= 100 {
			err = i.Batch(batch)
			if err != nil {
				return err
			}
			batch = i.NewBatch()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			fmt.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	// flush the last batch
	if batchCount > 0 {
		err = i.Batch(batch)
		if err != nil {
			return err
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	fmt.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}
