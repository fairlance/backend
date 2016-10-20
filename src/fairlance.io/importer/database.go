package importer

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fairlance.io/application"
	"github.com/jinzhu/gorm"
)

func reGenerateTestData(db gorm.DB, selectedType string) error {
	if err := deleteAllFromDB(db, selectedType); err != nil {
		return err
	}
	switch selectedType {
	case "jobs":
		// PostgreSQL only supports 65535 parameters
		for i := 0; i < 50; i++ {
			if err := db.Create(&application.Job{
				Name:      fmt.Sprintf("Job %d", i),
				Summary:   fmt.Sprintf("Job Summary %d", i),
				Details:   fmt.Sprintf("Job Description %d", i),
				ClientID:  1,
				Price:     123*i%200 + 200,
				StartDate: time.Now().Add(time.Duration(i*24+1) * time.Hour),
				Tags:      []string{fmt.Sprintf("tag_%d", i), fmt.Sprintf("tag_%d", i+i)},
			}).Error; err != nil {
				return err
			}
		}
	case "freelancers":
		for i := 0; i < 50; i++ {
			if err := db.Create(&application.Freelancer{
				User: application.User{
					FirstName: fmt.Sprintf("Name %d", i),
					LastName:  fmt.Sprintf("Last %d", i),
					Password:  fmt.Sprintf("Pass %d", i),
					Email:     fmt.Sprintf("email%d@mail.com", i),
				},
				HourlyRateFrom: 3,
				HourlyRateTo:   55,
				Timezone:       "UTC",
			}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteAllFromDB(db gorm.DB, selectedType string) error {
	switch selectedType {
	case "jobs":
		db.DropTableIfExists(&application.Job{})
		db.CreateTable(&application.Job{})
	case "freelancers":
		db.DropTableIfExists(&application.Freelancer{})
		db.CreateTable(&application.Freelancer{})
	}

	return nil
}

func doIndex(options Options, db gorm.DB, selectedType string) error {
	switch selectedType {
	case "jobs":
		jobsFromDB, _, err := getJobsFromDB(db, 0, -1)
		if err != nil {
			return err
		}

		if err = importDocs(options, "jobs", jobsFromDB); err != nil {
			return err
		}
	case "freelancers":
		freelancersFromDB, _, err := getFreelancersFromDB(db, 0, -1)
		if err != nil {
			return err
		}

		if err = importDocs(options, "freelancers", freelancersFromDB); err != nil {
			return err
		}
	}

	return nil
}

func getJobsFromDB(db gorm.DB, start, limit int) (map[string]interface{}, int, error) {
	jobs := []application.Job{}
	jobsMap := make(map[string]interface{})
	var total int
	dbQuery := db.Offset(start)
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if err := dbQuery.Find(&jobs).Error; err != nil {
		return jobsMap, total, err
	}

	for _, job := range jobs {
		id := strconv.FormatUint(uint64(job.ID), 10)
		jobsMap[id] = job
	}

	if err := db.Find(&[]application.Job{}).Count(&total).Error; err != nil {
		return jobsMap, total, err
	}

	return jobsMap, total, nil
}

func getFreelancersFromDB(db gorm.DB, start, limit int) (map[string]interface{}, int, error) {
	freelancersMap := make(map[string]interface{})
	freelancers := []application.Freelancer{}
	var total int
	dbQuery := db.Offset(start)
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if err := dbQuery.Find(&freelancers).Error; err != nil {
		return freelancersMap, total, err
	}

	for _, freelancer := range freelancers {
		id := strconv.FormatUint(uint64(freelancer.ID), 10)
		freelancersMap[id] = freelancer
	}

	if err := db.Find(&[]application.Freelancer{}).Count(&total).Error; err != nil {
		return freelancersMap, total, err
	}

	return freelancersMap, total, nil
}

func getDB(options Options) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"dbname=%s user=%s password=%s sslmode=disable", options.DBName, options.DBUser, options.DBPass))
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Opened DB")
	return db, nil
}
