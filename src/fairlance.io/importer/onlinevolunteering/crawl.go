package onlinevolunteering

import (
	"log"

	"fairlance.io/application"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

type item struct {
	Category string
	Link     string
	LinkText string
	Title    string
	Body     string
	Time     string
	Tags     []string
}

// GetJobs ...
func GetJobs() []application.Job {
	doc, err := goquery.NewDocument("https://www.onlinevolunteering.org/en/opportunities")
	if err != nil {
		log.Fatal(err)
	}
	jobs := []application.Job{}

	doc.Find(".opportunities-item").Each(func(i int, s *goquery.Selection) {
		category := s.Find(".category-head .name").Text()
		link, _ := s.Find("a.basic-link").Attr("href")
		linkText := s.Find("a.basic-link").Text()
		context := s.Find(".opportunity-content-wrapper")
		title := context.Find(".title h2").Text()
		body := context.Find(".body p").Text()
		description := s.Find(".description-block")
		time := description.Find(".time .number").Text()
		country := description.Find(".country").Text()
		tags := []string{strings.ToLower(strings.TrimSpace(country))}
		expertise := description.Find(".area-of-expertise").Text()
		expertise = strings.ToLower(strings.TrimSpace(expertise))
		expertiseSlice := strings.Split(expertise, " and ")
		if len(expertiseSlice) == 2 {
			tags = append(tags, strings.TrimSpace(expertiseSlice[1]))
		}
		for _, tag := range strings.Split(expertiseSlice[0], ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}

		job := application.Job{
			Name:    title,
			Summary: category + ", " + time + " hour/week",
			Details: body,
			Price:   0,
			Tags:    tags,
			Attachments: []application.Attachment{
				application.Attachment{
					Name: linkText,
					URL:  link,
				},
			},
		}

		jobs = append(jobs, job)
	})

	return jobs
}