package application

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start))
	})
}

func ContextAwareHandler(next http.Handler, appContext *ApplicationContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "context", appContext)
		next.ServeHTTP(w, r)
	})
}

func CORSHandler(next http.Handler, route Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			// todo: make configurable
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(route.AllowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		// Stop here for a Preflighted OPTIONS request.
		if r.Method == "OPTIONS" {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			WriteError(w, http.StatusBadRequest, jwt.ErrNoTokenInRequest)
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

		if err != nil || !token.Valid {
			WriteError(w, http.StatusBadRequest, errors.New("Not logged in."))
			return
		}
		next.ServeHTTP(w, r)
	})
}
