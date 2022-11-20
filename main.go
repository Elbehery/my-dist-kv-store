package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"playground/my-dist-kv-store/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", handlers.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", handlers.KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", handlers.KeyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9400", r))
}
