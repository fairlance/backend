package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"fairlance.io/application"
	"fairlance.io/messaging"
)

var (
	port       int
	secret     string
	projectURL string
)

func init() {
	flag.IntVar(&port, "port", 3005, "Specify the port to listen to.")
	flag.StringVar(&secret, "secret", "secret", "Secret string used for JWS.")
	flag.StringVar(&projectURL, "projectURL", "", "Project url to check for user access rights.")
	flag.Parse()

	f, err := os.OpenFile("/var/log/fairlance/messaging.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	router := mux.NewRouter()

	hub := messaging.NewHub(messaging.NewMessageDB())
	go hub.Run()

	router.Handle("/{name}/{room}/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		user := mux.Vars(r)["name"]
		messaging.NewServeWS(hub, &application.User{
			Model: application.Model{
				ID: uint(rand.Intn(100)),
			},
			FirstName: "User",
			LastName:  user,
		}, room).ServeHTTP(w, r)
	}))

	router.Handle("/{username}/{room}/send", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		username := mux.Vars(r)["username"]
		message := r.URL.Query().Get("message")
		hub.SendMessage(room, username, message)
	}))

	checkAccess := func(userID uint, room string) (bool, error) {
		response, err := http.Get(projectURL + "/" + room)
		if err != nil {
			return false, err
		}

		var responseStruct struct {
			Code    int                 `json:"code"`
			Data    application.Project `json:"data"`
			Success bool                `json:"success"`
		}
		err = json.NewDecoder(response.Body).Decode(&responseStruct)
		if err != nil || responseStruct.Code != http.StatusOK || responseStruct.Success == false {
			return false, err
		}

		found := false
		for _, user := range responseStruct.Data.Freelancers {
			if user.ID == userID {
				found = true
				break
			}
		}

		return found, nil
	}

	// requires a GET 'token' parameter
	router.Handle("/{room}/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messaging.NewWithTokenFromParams(func(token string) http.Handler {
			return application.NewAuthenticateWithClaims(secret, token, func(claims map[string]interface{}) http.Handler {
				return application.NewWithUserFromClaims(claims, func(user *application.User) http.Handler {
					return messaging.NewValidatedWithRoom(user, checkAccess, func(room string) http.Handler {
						return messaging.NewServeWS(hub, user, room)
					})
				})
			})
		}).ServeHTTP(w, r)
	}))

	// opts := &respond.Options{
	// 	Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
	// 		dataEnvelope := map[string]interface{}{"code": status}
	// 		if err, ok := data.(error); ok {
	// 			dataEnvelope["error"] = err.Error()
	// 			dataEnvelope["success"] = false
	// 		} else {
	// 			dataEnvelope["data"] = data
	// 			dataEnvelope["success"] = true
	// 		}
	// 		return status, dataEnvelope
	// 	},
	// }

	// // generates temporary token for certain room to be used for openening ws (not used)
	// // to be used with /{token}/ws
	// router.Handle("/{room}/token", opts.Handler(application.NewWithTokenFromHeader(func(token string) http.Handler {
	// 	return application.NewAuthenticateWithClaims(secret, token, func(claims map[string]interface{}) http.Handler {
	// 		return application.NewWithUserFromClaims(claims, func(user *application.User) http.Handler {
	// 			return messaging.NewValidatedWithRoom(user, func(room string) http.Handler {
	// 				return messaging.NewGenerateToken(secret, room, user)
	// 			})
	// 		})
	// 	})
	// })))

	// router.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	messaging.NewWithTokenFromParams(func(token string) http.Handler {
	// 		return application.NewAuthenticateWithClaims(secret, token, func(claims map[string]interface{}) http.Handler {
	// 			return messaging.NewWithRoomFromClaims(claims, func(room string) http.Handler {
	// 				return application.NewWithUserFromClaims(claims, func(user *application.User) http.Handler {
	// 					return messaging.NewServeWS(hub, user, room)
	// 				})
	// 			})
	// 		})
	// 	}).ServeHTTP(w, r)
	// }))
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))

}
