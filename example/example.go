package main

import (
	"log"
	"net/http"

	"gopkg.in/sauerbraten/crudapi.v2"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello there!\n"))
}

func main() {
	// storage
	storage := NewMapStorage()
	storage.AddMap("artists")
	storage.AddMap("albums")

	// create CRUD API routes
	api := crudapi.New(storage)

	// mount the API
	http.Handle("/api/", http.StripPrefix("/api", api))

	// mount a custom handler
	http.HandleFunc("/", hello)

	// start listening
	log.Println("server listening on localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
