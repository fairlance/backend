package application

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(appContext *ApplicationContext) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.Handler
		handler = context.ClearHandler(RecoverHandler(LoggerHandler(handler)))
		handler = CORSHandler(ContextAwareHandler(route.Handler, appContext), route)

		router.
			Methods(route.AllowedMethods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
