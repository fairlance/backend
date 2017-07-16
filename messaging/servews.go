package messaging

import (
	"errors"
	"log"
	"net/http"

	respond "gopkg.in/matryer/respond.v1"

	"github.com/fairlance/backend/middleware"
	"github.com/fairlance/backend/models"

	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appUser := context.Get(r, "user").(*models.User)
		projectID := context.Get(r, "projectID").(uint)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("could not upgrade connection: %v", err)
			return
		}
		newConnection := &userConn{
			userType:  appUser.Type,
			userID:    appUser.ID,
			conn:      conn,
			projectID: uint(projectID),
		}
		hub.register <- newConnection
	})
}

type userConn struct {
	userType  string
	userID    uint
	projectID uint
	conn      *websocket.Conn
}

func validateUser(hub *Hub) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			projectID := context.Get(r, "projectID").(uint)
			user := context.Get(r, "user").(*models.User)
			projectRoom, err := hub.getProjectRoom(projectID)
			if err != nil {
				log.Printf("could not get project room: %v", err)
				respond.With(w, r, http.StatusFailedDependency, errors.New("could not get project room"))
				return
			}
			if !projectRoom.isUserAllowed(user.Type, user.ID) {
				log.Printf("unauthorized user: %s %d", user.Type, user.ID)
				respond.With(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// type hasAccessFunc func(userID uint, userType, token, room string) (bool, error)

// func WithTokenFromParams(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token := r.URL.Query().Get("token")
// 		if token == "" {
// 			respond.With(w, r, http.StatusBadRequest, errors.New("valid token is missing from parameters"))
// 			return
// 		}
// 		context.Set(r, "token", token)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func WithRoom(hub *Hub) middleware.Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			vars := mux.Vars(r)
// 			roomID := vars["room"]
// 			if roomID == "" {
// 				respond.With(w, r, http.StatusBadRequest, "Room not provided.")
// 				return
// 			}
// 			if hub.rooms[roomID] == nil {
// 				room, err := hub.getARoom(roomID)
// 				if err != nil {
// 					log.Println(err)
// 					return
// 				}
// 				hub.rooms[roomID] = room
// 			}
// 			context.Set(r, "room", roomID)
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
