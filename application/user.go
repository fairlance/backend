package application

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/context"
	respond "gopkg.in/matryer/respond.v1"
)

func withUserToAdd(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var body struct {
			FirstName   string `json:"firstName" valid:"required"`
			LastName    string `json:"lastName" valid:"required"`
			Password    string `json:"password" valid:"required"`
			Email       string `json:"email" valid:"required,email"`
			Salutation  string `json:"salutation" valid:"required"`
			IsCompany   bool   `json:"isCompany"`
			CompanyName string `json:"companyName"`
		}
		if err := decoder.Decode(&body); err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		if ok, err := govalidator.ValidateStruct(body); ok == false || err != nil {
			errs := govalidator.ErrorsByField(err)
			respond.With(w, r, http.StatusBadRequest, errs)
			return
		}
		user := &User{
			FirstName:   body.FirstName,
			LastName:    body.LastName,
			Password:    body.Password,
			Email:       body.Email,
			Salutation:  body.Salutation,
			IsCompany:   body.IsCompany,
			CompanyName: body.CompanyName,
		}
		context.Set(r, "userToAdd", user)
		handler.ServeHTTP(w, r)
	})
}
