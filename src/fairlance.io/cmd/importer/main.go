package main

import (
	"flag"
	"fmt"
	"time"

	"fairlance.io/application"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/blevesearch/bleve"
)

var (
	dbName string
	dbUser string
	dbPass string
)

func main() {
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "fairlance", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "fairlance", "Db user's password.")
	flag.Parse()

	// open a new index
	jobsIndex, err := getIndex("jobs")
	if err != nil {
		panic(err)
	}

	freelancersIndex, err := getIndex("freelancers")
	if err != nil {
		panic(err)
	}

	db, err := getDB()
	if err != nil {
		panic(err)
	}

	jobsFromDB, err := jobs(db)
	if err != nil {
		panic(err)
	}
	for _, job := range jobsFromDB {
		fmt.Printf("Adding job %d ...\n", job.ID)
		jobsIndex.Index(string(job.ID), job)
	}

	freelancersFromDB, err := freelancers(db)
	if err != nil {
		panic(err)
	}
	for _, freelancer := range freelancersFromDB {
		fmt.Printf("Adding freelancer %d ...\n", freelancer.ID)
		freelancersIndex.Index(string(freelancer.ID), freelancer)
	}

	fmt.Println("Done.")
}

func jobs(db *gorm.DB) ([]application.Job, error) {

	db.DropTableIfExists(&application.Job{}, &application.Tag{})
	db.CreateTable(&application.Job{}, &application.Tag{})

	for i := 0; i < 100; i++ {
		db.Create(&application.Job{
			Name:        fmt.Sprintf("Job %d", i),
			Description: fmt.Sprintf("Job Description %d", i),
			ClientId:    1,
			Price:       123*i%200 + 200,
			StartDate:   time.Now().Add(time.Duration(i*24+1) * time.Hour),
			Tags: []application.Tag{
				application.Tag{Name: fmt.Sprintf("tag_%d", i)},
				application.Tag{Name: fmt.Sprintf("tag_%d", i+i)},
			},
		})
	}

	jobs := []application.Job{}
	if err := db.Preload("Tags").Preload("Client").Find(&jobs).Error; err != nil {
		return jobs, err
	}

	fmt.Printf("Found %d jobs ...\n", len(jobs))
	return jobs, nil
}

func freelancers(db *gorm.DB) ([]application.Freelancer, error) {
	freelancers := []application.Freelancer{}
	if err := db.Preload("Skills", "owner_type = ?", "freelancers").Find(&freelancers).Error; err != nil {
		return freelancers, err
	}

	fmt.Printf("Found %d freelancers ...\n", len(freelancers))
	return freelancers, nil
}

func getIndex(dbName string) (bleve.Index, error) {
	fmt.Printf("Opening %sIndex ...\n", dbName)
	index, err := bleve.Open("/tmp/" + dbName)
	if err != nil {
		fmt.Printf("%sIndex not found. Creating ...\n", dbName)
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New("/tmp/"+dbName, mapping)
		if err != nil {
			return index, err
		}
	}

	fmt.Printf("Opened %sIndex\n", dbName)
	return index, nil
}

func getDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", dbName, dbUser, dbPass))
	if err != nil {
		fmt.Println(err.Error())
		return db, err
	}

	fmt.Printf("Opened DB\n")
	return db, nil
}
