package application

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/matryer/respond.v1"
	"net/http"
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

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = opts.Handler(route.Handler)
		handler = CORSHandler(ContextAwareHandler(handler, appContext), route)
		handler = context.ClearHandler(LoggerHandler(handler))

		router.
			Methods(route.AllowedMethods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
