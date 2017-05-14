package middleware

import (
	"log"
	"net/http"
	"time"

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
