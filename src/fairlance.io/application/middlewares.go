package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
)

func IdHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if vars["id"] == "" {
			respond.With(w, r, http.StatusBadRequest, errors.New("Id not provided."))
			return
		}

		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		context.Set(r, "id", uint(id))
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
			w.Header().Set("Access-Control-Allow-Methods", route.Method+",OPTIONS")
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

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respond.With(w, r, http.StatusBadRequest, errors.New("Authorization header missing."))
			return
		}

		if tokenString[:7] != "Bearer " {
			respond.With(w, r, http.StatusBadRequest, errors.New("Authorization header must start with 'Bearer '."))
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

		if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			respond.With(w, r, http.StatusBadRequest, errors.New("Not logged in."))
			return
		} else {
			context.Set(r, "user", claims["user"])
		}

		next.ServeHTTP(w, r)
	})
}

//func HTTPAuthHandler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		user := ...
//		pass := ...
//		if !authenticated(w, r, user, pass) {
//			w.Header().Set("WWW-Authenticate", `Basic realm="FAIRLANCE"`)
//			w.WriteHeader(http.StatusUnauthorized)
//			w.Write([]byte("401 Unauthorized\n"))
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
//
//func authenticated(w http.ResponseWriter, r *http.Request, user string, pass string) bool {
//	authCredentials := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
//	if len(authCredentials) != 2 {
//		return false
//	}
//
//	credentials, err := base64.StdEncoding.DecodeString(authCredentials[1])
//	if err != nil {
//		return false
//	}
//
//	userAndPass := strings.SplitN(string(credentials), ":", 2)
//	if len(userAndPass) != 2 {
//		return false
//	}
//
//	return userAndPass[0] == user && userAndPass[1] == pass
//}
