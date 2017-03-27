package messaging

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"fairlance.io/application"

	respond "gopkg.in/matryer/respond.v1"

	jwt "github.com/dgrijalva/jwt-go"
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

// ServeWS builds websocket handler
func ServeWS(hub *Hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appUser := context.Get(r, "user").(*application.User)
		claims := context.Get(r, "claims").(jwt.MapClaims)
		userType := claims["userType"].(string)
		room := context.Get(r, "room").(string)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		newConnection := &userConn{
			id:   fmt.Sprintf("%s.%d", userType, appUser.ID),
			conn: conn,
			room: room,
			hub:  hub,
		}

		hub.register <- newConnection
	})
}

type hasAccessFunc func(userID uint, userType, token, room string) (bool, error)

// ValidateUser ...
func ValidateUser(hub *Hub, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roomName := context.Get(r, "room").(string)
		claims, ok := context.Get(r, "claims").(jwt.MapClaims)
		if !ok {
			log.Println("validate user: claims not of type jwt.MapClaims")
			respond.With(w, r, http.StatusInternalServerError, errors.New("could not validate user"))
			return
		}
		userType := claims["userType"].(string)
		user, ok := context.Get(r, "user").(*application.User)
		if !ok {
			log.Println("validate user: user not of type application.User")
			respond.With(w, r, http.StatusInternalServerError, errors.New("could not validate user"))
			return
		}

		room, ok := hub.rooms[roomName]
		if !ok {
			log.Println("room not found")
			respond.With(w, r, http.StatusNotFound, errors.New("room not found"))
			return
		}

		if !room.HasUser(user.ID, userType) {
			log.Println("unauthorized")
			respond.With(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// WithTokenFromParams ...
func WithTokenFromParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			respond.With(w, r, http.StatusBadRequest, errors.New("valid token is missing from parameters"))
			return
		}

		context.Set(r, "token", token)

		next.ServeHTTP(w, r)
	})
}

func WithRoom(hub *Hub, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		roomID := vars["room"]

		if roomID == "" {
			respond.With(w, r, http.StatusBadRequest, "Room not provided.")
			return
		}

		if hub.rooms[roomID] == nil {
			room, err := hub.getARoom(roomID)
			if err != nil {
				log.Println(err)
				return
			}
			hub.rooms[roomID] = room
		}

		context.Set(r, "room", roomID)

		next.ServeHTTP(w, r)
	})
}

//// NewGenerateToken ...
//func NewGenerateToken(secret string, room string, user *application.User) *GenerateToken {
//	return &GenerateToken{secret, room, user}
//}
//
//// GenerateToken ...
//type GenerateToken struct {
//	secret string
//	room   string
//	user   *application.User
//}
//
//func (t *GenerateToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	claims := make(map[string]interface{})
//	claims["user"] = *t.user
//	claims["room"] = t.room
//	token, err := application.CreateToken(claims, t.secret, time.Second*20)
//	if err != nil {
//		respond.With(w, r, http.StatusInternalServerError, err)
//		return
//	}
//
//	log.Println("new token for ws", token)
//
//	respond.With(w, r, http.StatusOK, struct {
//		Token string `json:"token"`
//	}{
//		Token: token,
//	})
//}
//
//// NewWithRoomFromClaims ...
//func NewWithRoomFromClaims(claims map[string]interface{}, next func(room string) http.Handler) *WithRoomFromClaims {
//	return &WithRoomFromClaims{claims, next}
//}
//
//// WithRoomFromClaims ...
//type WithRoomFromClaims struct {
//	claims map[string]interface{}
//	next   func(room string) http.Handler
//}
//
//func (wrfc *WithRoomFromClaims) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	room, ok := wrfc.claims["room"].(string)
//	if !ok || room == "" {
//		respond.With(w, r, http.StatusBadRequest, errors.New("valid room is missing from token"))
//		return
//	}
//
//	wrfc.next(room).ServeHTTP(w, r)
//}
