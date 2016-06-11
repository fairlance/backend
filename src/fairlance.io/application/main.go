package main

import (
	"flag"
	"net/http"
	"strconv"
)

var port int

func main() {
	flag.IntVar(&port, "port", 3001, "Specify the port to listen to.")
	flag.Parse()

	var appContext, err = NewContext("application")
	appContext.DropCreateFillTables()
	if err != nil {
		panic(err)
	}

	router := NewRouter(appContext)
	http.Handle("/", router)
	if err := http.ListenAndServe(":3001", nil); err != nil {
		panic(err)
	}
	panic(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
