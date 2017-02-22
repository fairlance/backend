package application

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

func withID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if vars["id"] == "" {
			respond.With(w, r, http.StatusBadRequest, "Id not provided.")
			return
		}

		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		context.Set(r, "id", uint(id))

		handler.ServeHTTP(w, r)
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

func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respond.With(w, r, http.StatusForbidden, errors.New("authorization header missing"))
			return
		}

		if tokenString[:7] != "Bearer " {
			respond.With(w, r, http.StatusForbidden, errors.New("authorization header must start with 'Bearer '"))
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
		context.Set(r, "userType", claims["userType"])

		next.ServeHTTP(w, r)
	})
}

// WithTokenFromHeader ...
func WithTokenFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respond.With(w, r, http.StatusForbidden, errors.New("authorization header missing"))
			return
		}

		if tokenString[:7] != "Bearer " {
			respond.With(w, r, http.StatusForbidden, errors.New("authorization header must start with 'Bearer '"))
			return
		}

		context.Set(r, "token", tokenString[7:])

		next.ServeHTTP(w, r)
	})
}

// AuthenticateTokenWithClaims ...
func AuthenticateTokenWithClaims(secret string, next http.Handler) http.Handler {
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

func getUserFomClaims(claims map[string]interface{}) (*User, error) {
	user := &User{}
	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		return user, errors.New("valid user is missing from token")
	}
	user.Model = Model{
		ID: uint(userMap["id"].(float64)),
	}
	user.Email = userMap["email"].(string)
	user.FirstName = userMap["firstName"].(string)
	user.LastName = userMap["lastName"].(string)
	createdAt, err := time.Parse(time.RFC3339, userMap["createdAt"].(string))
	if err != nil {
		return user, err
	}
	user.CreatedAt = createdAt

	updatedAt, err := time.Parse(time.RFC3339, userMap["updatedAt"].(string))
	if err != nil {
		return user, err
	}
	user.UpdatedAt = updatedAt

	return user, nil
}

// WithUserFromClaims ...
func WithUserFromClaims(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := context.Get(r, "claims").(jwt.MapClaims)
		user, err := getUserFomClaims(claims)
		if err != nil {
			respond.With(w, r, http.StatusBadRequest, err)
		}

		context.Set(r, "user", user)

		handler.ServeHTTP(w, r)
	})
}

// func HTTPAuthHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := "fairlance"
// 		pass := "fairlance"
// 		if !authenticated(w, r, user, pass) {
// 			w.Header().Set("WWW-Authenticate", `Basic realm="FAIRLANCE"`)
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("401 Unauthorized\n"))
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func authenticated(w http.ResponseWriter, r *http.Request, user string, pass string) bool {
// 	authCredentials := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
// 	if len(authCredentials) != 2 {
// 		return false
// 	}

// 	credentials, err := base64.StdEncoding.DecodeString(authCredentials[1])
// 	if err != nil {
// 		return false
// 	}

// 	userAndPass := strings.SplitN(string(credentials), ":", 2)
// 	if len(userAndPass) != 2 {
// 		return false
// 	}

// 	return userAndPass[0] == user && userAndPass[1] == pass
// }
