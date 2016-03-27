package main

import (
    app "fairlance.io/application"
    "net/http"
)


func main() {
    var appContext, err = app.NewContext("application")
    if err != nil {
        panic(err)
    }
    router := app.NewRouter(appContext)
    http.Handle("/", router)
    if err := http.ListenAndServe(":3001", nil); err != nil {
        panic(err)
    }
}
