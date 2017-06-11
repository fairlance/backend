package application

import (
	"github.com/fairlance/backend/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func NewRouter(appContext *ApplicationContext) *mux.Router {
	router := mux.NewRouter()
	publicRouter := router.PathPrefix("/public").Subrouter()
	for _, route := range routes {
		publicRouter.
			Methods([]string{route.Method, "OPTIONS"}...). // todo
			Path(route.Pattern).
			Name(route.Name).
			Handler(middleware.Chain(
				middleware.RecoverHandler,
				middleware.LoggerHandler,
				middleware.JSONEnvelope,
				middleware.CORSHandler,
				context.ClearHandler,
				contextAwareHandler(appContext),
			)(route.Handler))
	}
	privateRouter := router.PathPrefix("/private").Subrouter()
	for _, route := range privateRoutes {
		privateRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(middleware.Chain(
				middleware.RecoverHandler,
				middleware.LoggerHandler,
				middleware.JSONEnvelope,
				context.ClearHandler,
				contextAwareHandler(appContext),
			)(route.Handler))
	}

	return router
}
