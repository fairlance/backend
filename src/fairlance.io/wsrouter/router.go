package wsrouter

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var writeWait = 10 * time.Second
var readWait = 15 * time.Minute

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Examples:

// {"to":["freelancer.1"],"from":"freelancer.1","type":"notification","data":{"text":"hahahah"}}
// {"type":"read", "from":"freelancer.1", "to":["freelancer.1"], "data": {"timestamp":"1487627243358"}}
type Message struct {
	To        []string               `json:"to,omitempty"`
	From      string                 `json:"from,omitempty"`
	Type      string                 `json:"type,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
	Read      bool                   `json:"read"`
}

type User struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	send     chan Message
}

type Router struct {
	broadcast  chan Message
	register   chan User
	unregister chan User
	conf       RouterConf
}

func (r *Router) Handle(req *http.Request, conn *websocket.Conn) {
	usr := r.conf.CreateUser(req)
	usr.send = make(chan Message)
	r.register <- *usr
	go r.StartReading(*usr, conn)
	r.StartWriting(*usr, conn)
}

func NewRouter(conf RouterConf) *Router {
	return &Router{
		broadcast:  make(chan Message),
		register:   make(chan User),
		unregister: make(chan User),
		conf:       conf,
	}
}

type RouterConf struct {
	CreateUser   func(r *http.Request) *User
	Register     func(usr User) []Message
	Unregister   func(usr User)
	BuildMessage func(b []byte) *Message
	BroadcastTo  func(msg *Message) []User
}

// Run the Hub
func (r *Router) Run() {
	for {
		select {
		case usr := <-r.register:
			messages := r.conf.Register(usr)
			for _, message := range messages {
				usr.send <- message
			}
		case usr := <-r.unregister:
			r.conf.Unregister(usr)
		case msg := <-r.broadcast:
			users := r.conf.BroadcastTo(&msg)
			for _, user := range users {
				user.send <- msg
			}

		}
	}
}
func (r *Router) BroadcastMessage(msg Message) {
	r.broadcast <- msg
}
