package main

import (
	"log"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ramsgoli/uberextbackend/uber"
	"github.com/ramsgoli/uberextbackend/keys"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

// our main function
func main() {
	context := keys.Keys{UberKey: os.Getenv("UBERKEY")}

	router := mux.NewRouter()
	router.HandleFunc("/getEstimate", func(w http.ResponseWriter, r *http.Request) {uber.GetUberEstimate(w, r, &context)}).Methods("GET")
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), router))
}
