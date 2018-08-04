package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/ramsgoli/uberextbackend/getcoords"
    "github.com/ramsgoli/uberextbackend/uber"
)

// our main function
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/getCoord", getcoords.GetLocation).Methods("GET")
    router.HandleFunc("/getEstimate", uber.GetUberEstimate).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", router))
}
