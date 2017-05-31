package notification

import (
	"io/ioutil"
	"net/http"

	"github.com/fairlance/backend/middleware"
	"github.com/fairlance/backend/notification/wsrouter"
)

type Notifications struct {
	secret string
	router *wsrouter.Router
	users  map[string]wsrouter.User
	db     *mongoDB
}

func New(secret, mongoHost string) Notifications {
	users := make(map[string]wsrouter.User)
	db := newMongoDatabase(mongoHost, "notification")
	return Notifications{
		secret: secret,
		db:     db,
		router: wsrouter.NewRouter(newRouterConf(users, db)),
		users:  users,
	}
}

func (n Notifications) Run() {
	n.router.Run()
}

func (n Notifications) Handler() http.Handler {
	return middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.WithTokenFromParams,
		middleware.AuthenticateTokenWithClaims(n.secret),
		middleware.WithUserFromClaims,
	)(n.router.ServeWS())
}

func (n Notifications) SendHandler() http.Handler {
	return middleware.Chain(
		middleware.RecoverHandler,
		middleware.LoggerHandler,
		middleware.CORSHandler,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(err.Error()))
			}
			msg := buildMessage(b, n.users, n.db)
			n.router.BroadcastMessage(*msg)
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}))
}
