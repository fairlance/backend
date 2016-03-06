package application

import (
    "github.com/gorilla/mux"
    "net/http"
)

func NewRouter(appContext *ApplicationContext) *mux.Router {

    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
        handler = route.HandlerFunc
        handler = ContextAwareHandler(handler, appContext)
        handler = CORSHandler(handler)
        handler = LoggerHandler(handler, route.Name)

        router.
        Methods(getMethods(route)...).
        Path(route.Pattern).
        Name(route.Name).
        Handler(handler)
    }

    return router
}

func getMethods(route Route) []string {
    if route.Method != "GET" {
        // add OPTIONS for CORS
        return []string{route.Method, "OPTIONS"}
    }
    return []string{route.Method}
}
