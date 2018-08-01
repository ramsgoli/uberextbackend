package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/ramsgoli/uberextbackend/getcoords"
)

// our main function
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/getCoord", getcoords.GetLocation).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", router))
}
