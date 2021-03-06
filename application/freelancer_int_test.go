package application_test

// func TestIndexFreelancerWhenEmpty(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestIndexFreelancerWhenEmpty in short mode")
// 	}
// 	setUp()
// 	is := is.New(t)

// 	w := httptest.NewRecorder()
// 	r := getRequest("GET", "")
// 	app.IndexFreelancer(w, r)

// 	is.Equal(w.Code, http.StatusOK)
// 	var data []interface{}
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
// 	is.Equal(data, []interface{}{})
// }

// func TestIndexFreelancerWithFreelancers(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestIndexFreelancerWithFreelancers in short mode")
// 	}
// 	setUp()
// 	is := is.New(t)
// 	AddFreelancerToDB()
// 	AddFreelancerToDB()

// 	w := httptest.NewRecorder()
// 	r := getRequest("GET", "")
// 	app.IndexFreelancer(w, r)

// 	is.Equal(w.Code, http.StatusOK)
// 	var data []interface{}
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
// 	is.Equal(len(data), 2)
// }

// func TestAddFreelancer(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestAddFreelancer in short mode")
// 	}
// 	setUp()
// 	is := is.New(t)

// 	w := httptest.NewRecorder()
// 	r := getRequest("PUT", "")
// 	context.Set(r, "user", GetMockUser())

// 	app.AddFreelancer(w, r)

// 	is.Equal(w.Code, http.StatusOK)
// 	var data map[string]interface{}
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))

// 	user := data["user"].(map[string]interface{})
// 	is.Equal(data["type"], "freelancer")
// 	is.NotEqual(user["id"], 0)
// 	is.Equal(user["firstName"], "Pera")
// 	is.Equal(user["lastName"], "Peric")
// 	is.True(strings.HasSuffix(user["email"].(string), "pera@gmail.com"))

// 	freelancers := GetFreelancersFromDB()
// 	is.Equal(len(freelancers), 1)
// 	is.NotEqual(freelancers[0].ID, 0)
// 	is.Equal(freelancers[0].FirstName, "Pera")
// 	is.Equal(freelancers[0].LastName, "Peric")
// 	is.True(strings.HasSuffix(freelancers[0].Email, "pera@gmail.com"))
// }

// func TestDeleteFreelancer(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestDeleteFreelancer in short mode")
// 	}
// 	setUp()
// 	is := is.New(t)
// 	id := AddFreelancerToDB()

// 	w := httptest.NewRecorder()
// 	r := getRequest("POST", "")
// 	context.Set(r, "id", id)

// 	app.DeleteFreelancer(w, r)

// 	var data map[string]interface{}
// 	is.Equal(w.Code, http.StatusOK)
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))

// 	is.Equal(len(GetFreelancersFromDB()), 0)
// }

// func TestAddFreelancerUpdates(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping TestAddFreelancerUpdates in short mode")
// 	}
// 	setUp()
// 	is := is.New(t)
// 	id := AddFreelancerToDB()

// 	rBody := app.FreelancerUpdate{
// 		Skills: []app.Tag{
// 			app.Tag{Tag: "coolcat"},
// 			app.Tag{Tag: "pimp"},
// 		},
// 		Timezone:       "UTC",
// 		IsAvailable:    true,
// 		HourlyRateFrom: 2,
// 		HourlyRateTo:   20,
// 	}

// 	w := httptest.NewRecorder()
// 	r := getRequest("POST", "")
// 	context.Set(r, "id", id)
// 	context.Set(r, "updates", &rBody)

// 	app.AddFreelancerUpdates(w, r)

// 	freelancers := GetFreelancersFromDB()
// 	data := freelancers[0]

// 	is.Equal(data.Skills[0].Tag, "coolcat")
// 	is.Equal(data.Skills[1].Tag, "pimp")
// 	is.Equal(data.Timezone, "UTC")
// 	is.Equal(data.IsAvailable, true)
// 	is.Equal(data.HourlyRateFrom, 2)
// 	is.Equal(data.HourlyRateTo, 20)
// }

// func GetMockUser() *app.User {
// 	var email string
// 	rand.Seed(time.Now().UTC().UnixNano())
// 	email = strconv.Itoa(rand.Intn(100)) + "pera@gmail.com"
// 	return &app.User{
// 		FirstName: "Pera",
// 		LastName:  "Peric",
// 		Password:  "$2a$10$VJ8H9EYOIj9mnyW5mUm/nOWUrz/Rkak4/Ov3Lnw1GsAm4gmYU6sQu",
// 		Email:     email,
// 	}
// }

// func AddFreelancerToDB() uint {
// 	u := GetMockUser()
// 	f := &app.Freelancer{User: *u}
// 	appContext.FreelancerRepository.AddFreelancer(f)
// 	return f.ID
// }

// func GetFreelancersFromDB() []app.Freelancer {
// 	f, _ := appContext.FreelancerRepository.GetAllFreelancers()
// 	return f
// }
