package application

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
)

func NewRouter(appContext *ApplicationContext) *mux.Router {

	opts := &respond.Options{
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

	router := mux.NewRouter()
	for _, route := range routes {
		var handler http.Handler
		handler = opts.Handler(route.Handler)
		handler = CORSHandler(ContextAwareHandler(handler, appContext), route)
		handler = context.ClearHandler(handler)

		router.
			Methods([]string{route.Method, "OPTIONS"}...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
