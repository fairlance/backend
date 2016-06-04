package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) (*JobRepository, error) {
	repo := &JobRepository{db}

	return repo, nil
}

func (repo *JobRepository) GetAllJobs() ([]Job, error) {
	jobs := []Job{}
	repo.db.Find(&jobs)
	return jobs, nil
}
