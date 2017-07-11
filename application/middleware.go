package application

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fairlance/backend/middleware"
	"github.com/fairlance/backend/models"
	respond "gopkg.in/matryer/respond.v1"
)

func withUINT(param string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			if vars[param] == "" {
				respond.With(w, r, http.StatusBadRequest, fmt.Errorf("%s not provided", param))
				return
			}
			value, err := strconv.ParseUint(vars[param], 10, 32)
			if err != nil {
				respond.With(w, r, http.StatusBadRequest, err)
				return
			}
			context.Set(r, param, uint(value))
			next.ServeHTTP(w, r)
		})
	}
}

func contextAwareHandler(appContext *ApplicationContext) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Set(r, "context", appContext)
			next.ServeHTTP(w, r)
		})
	}
}

func whenFreelancer(next http.Handler) http.Handler {
	return middleware.WhenUserType("freelancer")(next)
}

func whenClient(next http.Handler) http.Handler {
	return middleware.WhenUserType("client")(next)
}

func whenBasedOnUserType(clientHandler, freelancerHandler middleware.Middleware) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := context.Get(r, "user").(*models.User)
			switch user.Type {
			case "client":
				clientHandler(next).ServeHTTP(w, r)
				return
			case "freelancer":
				freelancerHandler(next).ServeHTTP(w, r)
				return
			}
			err := fmt.Errorf("could not recognize user type: %s", user.Type)
			respond.With(w, r, http.StatusBadRequest, err)
		})
	}
}

func basedOnUserType(clientHandler http.Handler, freelancerHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user").(*models.User)
		switch user.Type {
		case "client":
			clientHandler.ServeHTTP(w, r)
			return
		case "freelancer":
			freelancerHandler.ServeHTTP(w, r)
			return
		}
		err := fmt.Errorf("could not recognize user type: %s", user.Type)
		respond.With(w, r, http.StatusBadRequest, err)
	})
}

func whenCurrentProjectStatus(status string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var project = context.Get(r, "project").(*Project)
			if project.Status != status {
				log.Printf("project status forbidden: wanted %s, found %s ", status, project.Status)
				respond.With(w, r, http.StatusNotFound, fmt.Errorf("project status forbidden: wanted %s, found %s ", status, project.Status))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func whenLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respond.With(w, r, http.StatusForbidden, errors.New("authorization header missing"))
			return
		}
		if tokenString[:7] != "Bearer " {
			respond.With(w, r, http.StatusForbidden, errors.New("authorization header invalid prefix"))
			return
		}
		var appContext = context.Get(r, "context").(*ApplicationContext)
		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(appContext.JwtSecret), nil
		})
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			respond.With(w, r, http.StatusUnauthorized, errors.New("not logged in"))
			return
		}
		user, err := getUserFomClaims(claims)
		if err != nil {
			log.Println("auth, get user from claims:", err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	})
}

func getUserFomClaims(claims map[string]interface{}) (*models.User, error) {
	user := &models.User{}
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		return user, errors.New("valid user is missing from token")
	}
	user.ID = uint(userMap["id"].(float64))
	user.Email = userMap["email"].(string)
	user.FirstName = userMap["firstName"].(string)
	user.LastName = userMap["lastName"].(string)
	user.Image = userMap["image"].(string)
	user.Type = claims["userType"].(string)
	// createdAt, err := time.Parse(time.RFC3339, userMap["createdAt"].(string))
	// if err != nil {
	// 	return user, err
	// }
	// user.CreatedAt = createdAt

	// updatedAt, err := time.Parse(time.RFC3339, userMap["updatedAt"].(string))
	// if err != nil {
	// 	return user, err
	// }
	// user.UpdatedAt = updatedAt

	return user, nil
}

func whenIDBelongsToUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = context.Get(r, "id").(uint)
		var user = context.Get(r, "user").(*models.User)
		if id != user.ID {
			respond.With(w, r, http.StatusForbidden, fmt.Errorf("user not allowed to access this entity"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}
