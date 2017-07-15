package messaging

import (
	"log"
	"net/http"

	"github.com/fairlance/backend/models"

	"strconv"

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
		projectRoom := context.Get(r, "room").(string)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("could not upgrade connection: %v", err)
			return
		}
		projectID, err := strconv.Atoi(projectRoom)
		if err != nil {
			log.Printf("could not get project id: %v", err)
			return
		}
		newConnection := &userConn{
			userType:  appUser.Type,
			userID:    appUser.ID,
			conn:      conn,
			projectID: uint(projectID),
			// hub:       hub,
		}
		hub.register <- newConnection
	})
}

type hasAccessFunc func(userID uint, userType, token, room string) (bool, error)

// func ValidateUser(hub *Hub) middleware.Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			roomName := context.Get(r, "room").(string)
// 			user, ok := context.Get(r, "user").(*models.User)
// 			if !ok {
// 				log.Println("validate user: user not of type application.User")
// 				respond.With(w, r, http.StatusInternalServerError, errors.New("could not validate user"))
// 				return
// 			}
// 			room, ok := hub.rooms[roomName]
// 			if !ok {
// 				log.Println("room not found")
// 				respond.With(w, r, http.StatusNotFound, errors.New("room not found"))
// 				return
// 			}
// 			if !room.HasUser(user) {
// 				log.Println("unauthorized")
// 				respond.With(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

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
