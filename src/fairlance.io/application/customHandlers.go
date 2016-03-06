package application

import (
    "net/http"
    "time"
    "log"
    "github.com/gorilla/context"
)

func LoggerHandler(handler http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        handler.ServeHTTP(w, r)
        log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
    })
}

func ContextAwareHandler(handler http.Handler, appContext *ApplicationContext) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        context.Set(r, "context", appContext)
        handler.ServeHTTP(w, r)
    })
}

func CORSHandler(handler http.Handler) http.Handler {
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
        // Lets Gorilla work
        handler.ServeHTTP(w, r)
    })
}
