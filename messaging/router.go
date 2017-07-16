package messaging

import (
	"encoding/json"
	"net/http"

	"github.com/fairlance/backend/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(hub *Hub, secret string) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/private/{projectID}/send", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(privateSend(hub)))
	router.Handle("/public/{projectID}/ws", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.WithUINT("projectID"),
		middleware.WithTokenFromParams,
		middleware.AuthenticateTokenWithUser(secret),
		validateUser(hub),
	)(ServeWS(hub)))
	return router
}

func privateSend(hub *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID := mux.Vars(r)["projectID"]
		var message Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("could not reade message from body"))
			return
		}
		r.Body.Close()
		hub.SendMessage(projectID, message)
	})
}
