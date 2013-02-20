package crudapi

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testing"
)

func TestAPI(t *testing.T) {

}

// Shows the usage of Storage and API.
func ExampleAPI() {
	// storage
	s := NewMapStorage()
	s.AddKind("a")
	s.AddKind("b")
	s.AddKind("l")
	s.AddKind("s")

	api := NewAPI(s)

	// routes
	r := mux.NewRouter()
	r.StrictSlash(true)

	/*
		POST creates,
		PUT updates,
		GET returns,
		DELETE deletes
	*/

	post := r.Methods("POST").Subrouter()
	get := r.Methods("GET").Subrouter()
	put := r.Methods("PUT").Subrouter()
	del := r.Methods("DELETE").Subrouter()

	// crud
	post.HandleFunc("/{kind}", api.Create)
	get.HandleFunc("/{kind}/{id}", api.Get)
	put.HandleFunc("/{kind}/{id}", api.Update)
	del.HandleFunc("/{kind}/{id}", api.Delete)

	// start listening
	log.Println("server listening on localhost:8080")
	http.ListenAndServe(":8080", r)
}
