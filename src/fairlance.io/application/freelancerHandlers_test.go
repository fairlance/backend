package application_test

import (
	"net/http/httptest"
	"testing"

	app "fairlance.io/application"
	"github.com/cheekybits/is"
	"github.com/gorilla/context"
)

func TestFreelancerReviewHandler(t *testing.T) {
	is := is.New(t)
	requestBody := `
	{
		"title":        "tetetetetet",
		"content":      "content",
		"rating":       2.4,
		"jobId":        2,
		"clientId":     2,
		"freelancerId": 12
	}`

	w := httptest.NewRecorder()
	r := getRequest("GET", requestBody)
	context.Set(r, "id", uint(12))
	app.FreelancerReviewHandler(emptyHandler).ServeHTTP(w, r)
	review := context.Get(r, "review").(*app.Review)
	is.Equal(review.ClientId, 2)
	is.Equal(review.Content, "content")
	is.Equal(review.Rating, 2.4)
	is.Equal(review.Title, "tetetetetet")
	is.Equal(review.FreelancerId, 12)
}

func TestFreelancerReferenceHandler(t *testing.T) {
	is := is.New(t)
	requestBody := `
	{
	  "title": "ttttt",
		"content": "ccccc",
		"media": {
			"image": "i",
			"video": "v"
		},
		"freelancerId": 13
	}`

	w := httptest.NewRecorder()
	r := getRequest("GET", requestBody)
	context.Set(r, "id", uint(13))
	app.FreelancerReferenceHandler(emptyHandler).ServeHTTP(w, r)
	reference := context.Get(r, "reference").(*app.Reference)
	is.Equal(reference.Title, "ttttt")
	is.Equal(reference.Content, "ccccc")
	is.Equal(reference.Media.Image, "i")
	is.Equal(reference.Media.Video, "v")
}

func TestFreelancerUpdateHandler(t *testing.T) {
	is := is.New(t)

	requestBody := `
	{
		"skills":        	[
			{
				"name": "coolcat"
			},
			{
				"name": "pimp"
			}
		],
		"timezone":      	"UTC",
		"isAvailable":      true,
		"hourlyRateFrom":   2,
		"hourlyRateTo":     20
	}`

	w := httptest.NewRecorder()
	r := getRequest("POST", requestBody)

	app.FreelancerUpdateHandler(emptyHandler).ServeHTTP(w, r)

	data := context.Get(r, "updates").(*app.FreelancerUpdate)

	is.Equal(data.Skills[0], app.Tag{Name: "coolcat"})
	is.Equal(data.Skills[1], app.Tag{Name: "pimp"})
	is.Equal(data.Timezone, "UTC")
	is.Equal(data.IsAvailable, true)
	is.Equal(data.HourlyRateFrom, 2)
	is.Equal(data.HourlyRateTo, 20)
}
