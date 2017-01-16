package messaging

import (
	"errors"
	"log"
	"net/http"

	"fairlance.io/application"

	respond "gopkg.in/matryer/respond.v1"

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
		room := context.Get(r, "room").(string)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		user := &user{
			hub:       hub,
			conn:      conn,
			send:      make(chan message, 256),
			username:  appUser.FirstName + " " + appUser.LastName,
			projectID: room,
			id:        appUser.ID,
		}

		hub.register <- user

		go user.startWriting()
		user.startReading()
	})
}

type hasAccessFunc func(userID uint, room string) (bool, error)

// ValidateUser ...
func ValidateUser(hasAccess hasAccessFunc, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := context.Get(r, "room").(string)
		user := context.Get(r, "user").(*application.User)

		// check if user can access the room
		ok, err := hasAccess(user.ID, room)
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, errors.New("could not check the room"))
			return
		}

		if !ok {
			log.Printf("room not allowed\n")
			respond.With(w, r, http.StatusUnauthorized, errors.New("room not allowed"))
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

func WithRoom(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if vars["room"] == "" {
			respond.With(w, r, http.StatusBadRequest, "Room not provided.")
			return
		}

		context.Set(r, "room", vars["room"])

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
