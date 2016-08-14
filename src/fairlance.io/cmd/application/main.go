package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"fairlance.io/application"
)

var port int
var dbName string
var dbUser string
var dbPass string
var secret string

func main() {
	flag.IntVar(&port, "port", 3001, "Specify the port to listen to.")
	flag.StringVar(&dbName, "dbName", "application", "DB name.")
	flag.StringVar(&dbUser, "dbUser", "", "DB user.")
	flag.StringVar(&dbPass, "dbPass", "", "Db user's password.")
	flag.StringVar(&secret, "secret", "secret", "Secret string used for JWS.")
	flag.Parse()

	if dbUser == "" || dbPass == "" {
		fmt.Println("dbUser or dbPass empty!")
		return
	}

	options := application.ContextOptions{dbName, dbUser, dbPass, secret}

	var appContext, err = application.NewContext(options)
	if err != nil {
		panic(err)
	}
	appContext.DropCreateFillTables()

	router := application.NewRouter(appContext)
	http.Handle("/", router)

	panic(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}