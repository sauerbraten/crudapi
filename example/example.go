package main

import (
	"github.com/gorilla/mux"
	"github.com/sauerbraten/crudapi"
	"log"
	"net/http"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello there!"))
}

func main() {
	// storage
	s := NewMapStorage()
	s.AddMap("artists")
	s.AddMap("albums")

	// guard
	g := MapGuard{map[string][]crudapi.Action{
		"artists": {crudapi.ActionCreate, crudapi.ActionGet, crudapi.ActionUpdate},
		"albums":  {crudapi.ActionCreate, crudapi.ActionGet, crudapi.ActionGetAll, crudapi.ActionUpdate},
	}}

	// router
	r := mux.NewRouter()

	// mounting the API
	crudapi.MountAPI(r.Host("localhost").Subrouter(), s, g)

	// custom handler
	r.HandleFunc("/", hello)

	// start listening
	log.Println("server listening on localhost:8080")
	log.Println("API on api.localhost:8080/v1/")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}
