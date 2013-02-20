package crudapi

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testing"
)

func TestAPI(t *testing.T) {

}

// Put this code into a main.go, fix imports and stuff.
// When the server is running, try the following commands
//
// curl -i -X POST -d '{"id":"gorillaz","resource":{"name":"Gorillaz","albums":["the-fall"]}}' http://localhost:8080/artist
//
// curl -i -X POST -d '{"id":"plastic-beach","resource":{"title":"Plastic Beach","by":"gorillaz","songs":["on-melancholy-hill","stylo"]}}' http://localhost:8080/artist
//
// curl -i -X GET http://localhost:8080/artist/gorillaz
//
// curl -i -X PUT -d '{"name":"Gorillaz","albums":["plastic-beach","the-fall"]}' http://localhost:8080/artist/gorillaz
//
// Note the returned HTTP codes: '201 Created' when POSTing, '200 OK' when GETting and PUTting.
// There is also '404 Not Found' if either the kind of data you are posting (for example 'artist' and 'album' in the URLs) is unkown or there is no resource with the specified id ('gorillaz' in the GET request). In that case a JSON object containing an "error" field is returned, i.e.: {"error":"resource not found"} or {"error":"kind not found"}.
// '400 Bad Request' is returned when either the POSTed or PUTted JSON data is malformed and cannot be parsed or when you are POSTing/PUTting without an "id" field in the top-level JSON object.
// '409 Conflict' and {"error":"resource already exists"} as response means, well, that you POSTed a resource with an "id" that is already in use.
//
// Server responses are always a JSON object, containing one or more of the following fields:
// "error": specifies the error that occured, if any
// "id": the ID of the newly created or updated resource
// "resource": the requested resource (used when GETting resources)
//
func ExampleAPI() {
	// storage
	s := NewMapStorage()
	s.AddKind("artist")
	s.AddKind("album")

	api := NewAPI(s)

	// routes
	r := mux.NewRouter()
	r.StrictSlash(true)

	/*
		POST creates,
		GET returns,
		PUT updates,
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
