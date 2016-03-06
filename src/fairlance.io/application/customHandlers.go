package application

import (
    "net/http"
    "time"
    "log"
    "github.com/gorilla/context"
    "github.com/dgrijalva/jwt-go"
    "fmt"
)

func LoggerHandler(next http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
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
            w.Header().Set("Access-Control-Allow-Origin", "*")
            // todo: make this list configurable per route
            w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
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
            w.WriteHeader(http.StatusForbidden)
            w.Write([]byte(jwt.ErrNoTokenInRequest.Error()))
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
            w.WriteHeader(http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}