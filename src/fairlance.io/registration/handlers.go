package registration

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
)

type AppHandler struct {
	Context *RegistrationContext
	Handle  func(*RegistrationContext, http.ResponseWriter, *http.Request) error
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: make this configurable
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := ah.Handle(ah.Context, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.Context.Logger.Println(err)
	}
}

func IndexHandler(context *RegistrationContext, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RegistrationError{"Method not allowed! Use GET"})
		return nil
	}

	users, err := context.RegisteredUserRepository.GetAllRegisteredUsers()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
	return nil
}

func RegisterHandler(context *RegistrationContext, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RegistrationError{"Method not allowed! Use POST"})
		return nil
	}

	email, err := getEmailFromRequest(w, r)
	if err != nil {
		return err
	}

	if email != "" {

		// validate email first
		if !govalidator.IsEmail(email) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RegistrationError{"Email not valid!"})
			return nil
		}

		registerTime := time.Now()
		err := context.RegisteredUserRepository.AddRegisteredUser(RegisteredUser{email, registerTime})
		if err != nil {
			if mgo.IsDup(err) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(RegistrationError{"Email exists!"})
				return nil
			}

			return err
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(RegisteredUser{email, registerTime})
		context.Mailer.SendWelcomeMessage(email)
		return nil
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(RegistrationError{"Email missing!"})
	return nil
}

func getEmailFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Header.Get("Content-Type") == "application/json" {
		defer r.Body.Close()
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RegistrationError{"Request not valid JSON!"})
			return "", err
		}
		return data["email"], nil
	}

	return r.FormValue("email"), nil
}

func Authenticated(w http.ResponseWriter, r *http.Request, user string, pass string) bool {
	authCredentials := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authCredentials) != 2 {
		return false
	}

	credentials, err := base64.StdEncoding.DecodeString(authCredentials[1])
	if err != nil {
		return false
	}

	userAndPass := strings.SplitN(string(credentials), ":", 2)
	if len(userAndPass) != 2 {
		return false
	}

	return userAndPass[0] == user && userAndPass[1] == pass
}
