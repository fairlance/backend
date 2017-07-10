package middleware

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fairlance/backend/models"
	"github.com/gorilla/context"

	jwt "github.com/dgrijalva/jwt-go"
	respond "gopkg.in/matryer/respond.v1"
)

var opts = &respond.Options{
	Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
		dataEnvelope := map[string]interface{}{"code": status}
		if err, ok := data.(error); ok {
			dataEnvelope["error"] = err.Error()
			dataEnvelope["success"] = false
		} else {
			dataEnvelope["data"] = data
			dataEnvelope["success"] = true
		}
		return status, dataEnvelope
	},
}

type Middleware func(http.Handler) http.Handler

func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func JSONEnvelope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts.Handler(next).ServeHTTP(w, r)
	})
}

func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			// todo: make configurable
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here for a Preflighted OPTIONS request.
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start))
	})
}

func RecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				respond.With(w, r, http.StatusInternalServerError, nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func HTTPAuthHandler(user, password string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !authenticated(r, user, password) {
				w.Header().Set("WWW-Authenticate", `Basic realm="FAIRLANCE"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 Unauthorized"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func authenticated(r *http.Request, user, password string) bool {
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
	return userAndPass[0] == user && userAndPass[1] == password
}

func WithTokenFromHeader(next http.Handler) http.Handler {
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
		context.Set(r, "token", tokenString[7:])
		next.ServeHTTP(w, r)
	})
}

func WithTokenFromParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			respond.With(w, r, http.StatusBadRequest, errors.New("valid token is missing from parameters"))
			return
		}
		context.Set(r, "token", token)
		next.ServeHTTP(w, r)
	})
}

func AuthenticateTokenWithClaims(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenFromContext := context.Get(r, "token").(string)
			token, err := jwt.Parse(tokenFromContext, func(jwtToken *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
				}
				return []byte(secret), nil
			})
			if err != nil {
				respond.With(w, r, http.StatusUnauthorized, err)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				respond.With(w, r, http.StatusUnauthorized, errors.New("not logged in"))
				return
			}
			context.Set(r, "claims", claims)
			next.ServeHTTP(w, r)
		})
	}
}

func WithUserFromClaims(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &models.User{}
		claims, ok := context.Get(r, "claims").(jwt.MapClaims)
		if !ok {
			respond.With(w, r, http.StatusInternalServerError, errors.New("claims are missing from token"))
			return
		}
		userMap, ok := claims["user"].(map[string]interface{})
		if !ok {
			respond.With(w, r, http.StatusInternalServerError, errors.New("valid user is missing from token"))
			return
		}
		user.ID = uint(userMap["id"].(float64))
		user.Email = userMap["email"].(string)
		user.FirstName = userMap["firstName"].(string)
		user.LastName = userMap["lastName"].(string)
		user.Type = claims["userType"].(string)
		context.Set(r, "user", user)
		handler.ServeHTTP(w, r)
	})
}

func WhenUserType(allowedUserType string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := context.Get(r, "user").(*models.User)
			if user.Type != allowedUserType {
				respond.With(w, r, http.StatusForbidden, errors.New("unauthorized access"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func HTTPMethod(method string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				respond.With(w, r, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
