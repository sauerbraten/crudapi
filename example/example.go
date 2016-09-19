package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sauerbraten/crudapi"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello there!"))
}

func main() {
	// storage
	storage := NewMapStorage()
	storage.AddMap("artists")
	storage.AddMap("albums")

	// router
	r := mux.NewRouter()

	// mounting the API
	crudapi.Mount(r.Host("localhost").Subrouter(), storage, nil)

	// custom handler
	r.HandleFunc("/", hello)

	// start listening
	log.Println("server listening on localhost:8080")
	log.Println("API on localhost:8080/")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}
