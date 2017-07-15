package messaging

import (
	"encoding/json"
	"net/http"

	"github.com/fairlance/backend/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(hub *Hub, secret string) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/private/{room}/send", middleware.Chain(
		middleware.RecoverHandler,
		middleware.HTTPMethod(http.MethodPost),
	)(privateSend(hub)))
	// requires a GET 'token' parameter
	router.Handle("/public/{room}/ws", middleware.Chain(
		middleware.RecoverHandler,
		// messaging.WithRoom(hub),
		// messaging.WithTokenFromParams,
		middleware.WithTokenFromHeader,
		middleware.AuthenticateTokenWithClaims(secret),
		middleware.WithUserFromClaims,
		// messaging.ValidateUser(hub),
	)(ServeWS(hub)))
	return router
}

func privateSend(hub *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		var message Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("could not reade message from body"))
			return
		}
		r.Body.Close()
		hub.SendMessage(room, message)
	})
}
