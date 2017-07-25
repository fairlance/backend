package notification

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fairlance/backend/middleware"
	"github.com/fairlance/backend/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRouter(hub *Hub, secret string) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/public/ws", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.WithTokenFromParams,
		middleware.AuthenticateTokenWithUser(secret),
	)(addUserHandler(hub)))
	router.Handle("/public/", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.WithTokenFromParams,
		middleware.AuthenticateTokenWithUser(secret),
	)(addUserHandler(hub)))
	router.Handle("/private/send", middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.CORSHandler,
	)(sendMessageHandler(hub)))
	return router
}

func addUserHandler(hub *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		user := context.Get(r, "user").(*models.User)
		hub.addUser(user, conn)
	})
}

func sendMessageHandler(hub *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(err.Error()))
			}
			hub.sendMessage(b)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	})
}
