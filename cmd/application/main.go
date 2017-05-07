package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fairlance/backend/application"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	port            int
	dbName          string
	dbUser          string
	dbPass          string
	secret          string
	notificationURL string
	messagingURL    string
	searcherURL     string
)

func init() {
	f, err := os.OpenFile("/var/log/fairlance/application.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)

	flag.IntVar(&port, "port", 3001, "Specify the port to listen to.")
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "", "Db user's password.")
	flag.StringVar(&secret, "secret", "secret", "Secret string used for JWS.")
	flag.StringVar(&notificationURL, "notificationUrl", "localhost:3007", "Notification endpoint.")
	flag.StringVar(&messagingURL, "messagingUrl", "localhost:3007", "Messaging endpoint.")
	flag.StringVar(&searcherURL, "searcherUrl", "localhost:3003", "Url of the searcher.")
	flag.Parse()

	if dbUser == "" || dbPass == "" {
		log.Fatalln("dbUser or dbPass empty!")
	}
}

func main() {
	options := application.ContextOptions{
		DbName:          dbName,
		DbUser:          dbUser,
		DbPass:          dbPass,
		Secret:          secret,
		NotificationURL: notificationURL,
		MessagingURL:    messagingURL,
		SearcherURL:     searcherURL,
	}

	var appContext, err = application.NewContext(options)
	if err != nil {
		panic(err)
	}
	appContext.DropCreateFillTables()

	router := application.NewRouter(appContext)
	http.Handle("/", router)

	panic(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
