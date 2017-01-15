package messaging

import (
	"errors"
	"log"
	"net/http"

	respond "gopkg.in/matryer/respond.v1"

	"fairlance.io/application"

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

// NewServeWS creates new websocket handler
func NewServeWS(hub *Hub, user *application.User, room string) *ServeWS {
	return &ServeWS{hub, user, room}
}

// ServeWS has websocket handler
type ServeWS struct {
	hub  *Hub
	user *application.User
	room string
}

func (sws *ServeWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	user := &user{
		hub:       sws.hub,
		conn:      conn,
		send:      make(chan message, 256),
		username:  sws.user.FirstName + " " + sws.user.LastName,
		projectID: sws.room,
		id:        sws.user.ID,
	}

	sws.hub.register <- user

	go user.startWriting()
	user.startReading()
}

type hasAccessFunc func(userID uint, room string) (bool, error)

// NewValidatedWithRoom ...
func NewValidatedWithRoom(user *application.User, hasAccess hasAccessFunc, next func(room string) http.Handler) *ValidatedWithRoom {
	return &ValidatedWithRoom{user, hasAccess, next}
}

// ValidatedWithRoom provides room name from parameters
type ValidatedWithRoom struct {
	user      *application.User
	hasAccess hasAccessFunc
	next      func(room string) http.Handler
}

func (vwr *ValidatedWithRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["room"]

	// check if user can access the room
	ok, err := vwr.hasAccess(vwr.user.ID, roomID)
	if err != nil {
		log.Println(err)
		respond.With(w, r, http.StatusInternalServerError, errors.New("could not check the room"))
		return
	}

	if !ok {
		respond.With(w, r, http.StatusUnauthorized, errors.New("room not allowed"))
		return
	}

	vwr.next(roomID).ServeHTTP(w, r)
}

// NewWithTokenFromParams ...
func NewWithTokenFromParams(next func(token string) http.Handler) *WithTokenFromParams {
	return &WithTokenFromParams{next}
}

// WithTokenFromParams ...
type WithTokenFromParams struct {
	next func(token string) http.Handler
}

func (wrfc *WithTokenFromParams) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respond.With(w, r, http.StatusBadRequest, errors.New("valid token is missing from parameters"))
		return
	}

	wrfc.next(token).ServeHTTP(w, r)
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
