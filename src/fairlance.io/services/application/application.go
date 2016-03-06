package main

import (
    app "fairlance.io/application"
    "net/http"
)


func main() {
    var appContext = app.NewContext("application")
    router := app.NewRouter(appContext)
    http.Handle("/", router)
    if err := http.ListenAndServe(":3001", nil); err != nil {
        panic(err)
    }
}
