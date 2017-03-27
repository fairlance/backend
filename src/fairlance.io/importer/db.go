package importer

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"fairlance.io/application"
	"fairlance.io/importer/onlinevolunteering"
	"github.com/jinzhu/gorm"
)

func reGenerateTestData(db gorm.DB, selectedType string) error {
	if err := deleteAllFromDB(db, selectedType); err != nil {
		return err
	}
	switch selectedType {
	case "jobs":
		clients := []application.Client{}
		for i := 0; i < 5; i++ {
			clients = append(clients, application.Client{
				User: application.User{
					FirstName: fmt.Sprintf("Clint %d", i),
					LastName:  fmt.Sprintf("Clientwood %d", i),
					Email:     fmt.Sprintf("email%d@email.com", i),
				},
				Rating: float64(i % 5),
			})
			if err := db.Create(&clients[i]).Error; err != nil {
				return err
			}
		}

		jobs := onlinevolunteering.GetJobs()
		for i := 0; i < len(jobs); i++ {
			job := &jobs[i]
			job.ClientID = uint(i%5 + 1)
			job.StartDate = time.Now().Add(time.Duration(i*24+1) * time.Hour)
			job.Deadline = time.Now().Add(time.Duration(i*24*i+1) * time.Hour)
			if err := db.Create(job).Error; err != nil {
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
		db.DropTableIfExists(&application.Job{}, &application.Client{}, &application.Example{}, &application.Attachment{})
		db.CreateTable(&application.Job{}, &application.Client{}, &application.Example{}, &application.Attachment{})
	case "freelancers":
		db.DropTableIfExists(&application.Freelancer{})
		db.CreateTable(&application.Freelancer{})
	}

	return nil
}

func doImport(options Options, db gorm.DB, selectedType string) error {
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
	if err := dbQuery.Order("id ASC").Preload("Client").Find(&jobs).Error; err != nil {
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
	if err := dbQuery.Order("id ASC").Find(&freelancers).Error; err != nil {
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

func getDocFromDB(db gorm.DB, docType, docID string) (map[string]interface{}, error) {
	docMap := make(map[string]interface{})
	switch docType {
	case "jobs":
		entity := &application.Job{}
		if err := db.Find(entity, docID).Error; err != nil {
			return nil, err
		}
		docMap["ID"] = entity.ID
		docMap["Name"] = entity.Name
		docMap["ClientID"] = entity.ClientID
		docMap["CreatedAt"] = entity.CreatedAt
		docMap["UpdatedAt"] = entity.UpdatedAt
		docMap["Details"] = entity.Details
		docMap["IsActive"] = entity.IsActive
		docMap["Price"] = entity.Price
		docMap["StartDate"] = entity.StartDate
		docMap["Tags"] = entity.Tags
		docMap["Summary"] = entity.Summary
		docMap["Number of applications"] = len(entity.JobApplications)

		return docMap, nil
	case "freelancers":
		entity := &application.Freelancer{}
		if err := db.Find(entity, docID).Error; err != nil {
			return nil, err
		}

		docMap["ID"] = entity.ID
		docMap["FirstName"] = entity.FirstName
		docMap["LastName"] = entity.LastName
		docMap["IsAvailable"] = entity.IsAvailable
		docMap["CreatedAt"] = entity.CreatedAt
		docMap["Email"] = entity.Email
		docMap["HourlyRateFrom"] = entity.HourlyRateFrom
		docMap["HourlyRateTo"] = entity.HourlyRateTo
		docMap["Rating"] = entity.Rating
		docMap["UpdatedAt"] = entity.UpdatedAt
		docMap["Number of applications"] = len(entity.JobApplications)
		docMap["Number of projects"] = len(entity.Projects)

		return docMap, nil
	}

	return nil, errors.New("unknown type " + docType)
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
