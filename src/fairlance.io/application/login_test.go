package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	emptyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

var badBodyTest = []struct {
	in string
}{
	{"{}"},
	{""},
	{"{bad json}"},
	{`{"email":"no@password.com"}`},
	{`{"password":"noemail"}`},
}

func TestLoginWithBadEmailAndPassword(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &userRepositoryMock{}
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	for _, data := range badBodyTest {
		t.Run(data.in, func(t *testing.T) {
			r := getRequest(userContext, data.in)
			w := httptest.NewRecorder()

			Login(w, r)

			is.Equal(w.Code, http.StatusBadRequest)
		})
	}
}

func TestLoginWenUnauthorized(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &userRepositoryMock{}
	userRepositoryMock.CheckCredentialsCall.Returns.Error = errors.New("unauthorized")
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	r := getRequest(userContext, `
		{
			"email": "not@important.com",
			"password": "notimportant"
		}
	`)
	w := httptest.NewRecorder()

	Login(w, r)

	is.Equal(w.Code, http.StatusUnauthorized)
}

func TestLoginCallsCheckCredentialsWithCorrectEmailAndPassword(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &userRepositoryMock{}
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	r := getRequest(userContext, `
		{
			"email": "email@mail.com",
			"password": "password"
		}
	`)
	w := httptest.NewRecorder()

	Login(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.Equal(userRepositoryMock.CheckCredentialsCall.Receives.Email, "email@mail.com")
	is.Equal(userRepositoryMock.CheckCredentialsCall.Receives.Password, "password")
}

func TestLogin(t *testing.T) {
	is := is.New(t)
	userRepositoryMock := &userRepositoryMock{}
	userRepositoryMock.CheckCredentialsCall.Returns.UserType = "freelancer"
	userRepositoryMock.CheckCredentialsCall.Returns.User = User{
		Model: Model{
			ID: 1,
		},
		FirstName: "firstname",
		LastName:  "lastname",
		Password:  "password",
		Email:     "email@mail.com",
	}
	userContext := &ApplicationContext{
		UserRepository: userRepositoryMock,
	}

	r := getRequest(userContext, `
		{
			"email": "not@important.com",
			"password": "notimportant"
		}
	`)
	w := httptest.NewRecorder()

	Login(w, r)

	is.Equal(w.Code, http.StatusOK)
	var body map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
	is.Equal(body["id"], 1)
	is.OK(body["token"])
	is.Equal(body["type"], "freelancer")
}

// func TestIdHandler(t *testing.T) {
// 	is := is.New(t)

// 	r, _ := http.NewRequest("GET", "/1", nil)
// 	w := httptest.NewRecorder()
// 	router := mux.NewRouter()
// 	router.Handle("/{id}", app.IdHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := context.Get(r, "id").(uint)
// 		is.Equal(id, 1)
// 	}))).Methods("GET")
// 	router.ServeHTTP(w, r)
// }

// func TestUserHandler(t *testing.T) {
// 	is := is.New(t)
// 	requestBody := `
// 	{
// 	  "password": "123",
// 	  "email": "pera@gmail.com",
// 	  "firstName":"Pera",
// 	  "lastName":"Peric"
// 	}`

// 	w := httptest.NewRecorder()
// 	r := getRequest("GET", requestBody)
// 	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
// 	app.RegisterUserHandler(emptyHandler).ServeHTTP(w, r)
// 	user := context.Get(r, "user").(*app.User)
// 	is.Equal(user.FirstName, "Pera")
// 	is.Equal(user.LastName, "Peric")
// 	is.Equal(user.Email, "pera@gmail.com")
// }

// func TestUserHandlerWithInvalidBody(t *testing.T) {
// 	is := is.New(t)
// 	requestBody := `
// 	{
// 		"empty": "invalid body"
// 	}`

// 	w := httptest.NewRecorder()
// 	r := getRequest("GET", requestBody)
// 	app.RegisterUserHandler(emptyHandler).ServeHTTP(w, r)
// 	is.Equal(w.Code, http.StatusBadRequest)
// 	var errorBody map[string]string
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &errorBody))
// 	is.OK(errorBody["Email"])
// 	is.OK(errorBody["FirstName"])
// 	is.OK(errorBody["LastName"])
// 	is.OK(errorBody["Password"])
// }

// func TestUserHandlerWithInvalidEmail(t *testing.T) {
// 	is := is.New(t)
// 	requestBody := `
// 	{
// 	  "email": "invalid email",
// 	  "password": "123",
// 	  "firstName":"Pera",
// 	  "lastName":"Peric"
// 	}`

// 	w := httptest.NewRecorder()
// 	r := getRequest("GET", requestBody)
// 	app.RegisterUserHandler(emptyHandler).ServeHTTP(w, r)
// 	is.Equal(w.Code, http.StatusBadRequest)
// 	var body map[string]string
// 	is.NoErr(json.Unmarshal(w.Body.Bytes(), &body))
// 	is.OK(body["Email"])
// }

// func buildTestContext() *app.ApplicationContext {
// 	options := app.ContextOptions{"application_test", "fairlance", "fairlance", "secret"}
// 	context, err := app.NewContext(options)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return context
// }

// func setUp() {
// 	appContext = buildTestContext()
// 	appContext.DropTables()
// 	appContext.CreateTables()
// }
