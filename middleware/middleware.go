package middleware

import (
	"net/http"

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
