// Package crudapi implements HTTP handlers a minimalistic RESTful API offering Create, Read, Update, and Delete (→CRUD) handlers.

/*
See http://en.wikipedia.org/wiki/Create,_read,_update_and_delete for more information.

Note: Read is called Get in this package, but CGUD is hard to pronounce.

Example

Put this code into a 'main.go' file:

	package main

	import (
		"github.com/gorilla/mux"
		"github.com/sauerbraten/crudapi"
		"log"
		"net/http"
	)

	func main() {
		// storage
		s := crudapi.NewMapStorage()
		s.AddKind("artist")
		s.AddKind("album")

		api := crudapi.NewAPI(s)

		// routes
		r := mux.NewRouter()
		r.StrictSlash(true)

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

When the server is running, try the following commands

	curl -i -X POST -d '{"id":"gorillaz","resource":{"name":"Gorillaz","albums":["the-fall"]}}' http://localhost:8080/artist

	curl -i -X POST -d '{"id":"plastic-beach","resource":{"title":"Plastic Beach","by":"gorillaz","songs":["on-melancholy-hill","stylo"]}}' http://localhost:8080/artist

	curl -i -X GET http://localhost:8080/artist/gorillaz

	curl -i -X PUT -d '{"name":"Gorillaz","albums":["plastic-beach","the-fall"]}' http://localhost:8080/artist/gorillaz

Note the returned HTTP codes: '201 Created' when POSTing, '200 OK' when GETting and PUTting.
There is also '404 Not Found' if either the kind of data you are posting (for example 'artist' and 'album' in the URLs) is unkown or there is no resource with the specified id ('gorillaz' in the GET request). In that case a JSON object containing an "error" field is returned, i.e.: {"error":"resource not found"} or {"error":"kind not found"}.
'400 Bad Request' is returned when either the POSTed or PUTted JSON data is malformed and cannot be parsed or when you are POSTing/PUTting without an "id" field in the top-level JSON object.
'409 Conflict' and {"error":"resource already exists"} as response means, well, that you POSTed a resource with an "id" that is already in use.

Server responses are always a JSON object, containing one or more of the following fields:

	"error":     specifies the error that occured, if any
	"id":        the ID of the newly created or updated resource
	"resource":  the requested resource (used when GETting resources)
*/
package crudapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sauerbraten/slugify"
	"log"
	"net/http"
)

type apiResponse struct {
	Error    string      `json:"error,omitempty"`
	Id       string      `json:"id,omitempty"`
	Resource interface{} `json:"resource,omitempty"`
}

// API exposes the CRUD handlers.
type API struct {
	s Storage // the API's storage

	// Create is meant to handle POST requests. It returns '400 Bad Request', '404 Not Found', '409 Conflict' or '201 Created'.
	Create func(resp http.ResponseWriter, req *http.Request)

	// Get is meant to retrieve resources (HTTP GET). It returns '404 Not Found' or '200 OK'.
	Get func(resp http.ResponseWriter, req *http.Request)

	// Update is meant to handle PUTs. It returns '400 Bad Request', '404 Not Found' or '200 OK'.
	Update func(resp http.ResponseWriter, req *http.Request)

	// Delete is for handling DELETE requests. Possible HTTP status codes are '404 Not Found' and '200 OK'.
	Delete func(resp http.ResponseWriter, req *http.Request)
}

// Returns an API relying on the given Storage.
func NewAPI(s Storage) *API {
	api := &API{}

	api.s = s

	api.Create = func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		kind := vars["kind"]
		enc := json.NewEncoder(resp)
		dec := json.NewDecoder(req.Body)

		// read body and parse into interface{}
		var payload map[string]interface{}
		err := dec.Decode(&payload)

		if err != nil {
			log.Println(err)
			resp.WriteHeader(400) // Bad Request
			enc.Encode(apiResponse{"malformed json", "", nil})
			return
		}

		// create ID
		resourceName, ok := payload["id"]
		if !ok {
			resp.WriteHeader(400) // Bad Request
			enc.Encode(apiResponse{"no id field given", "", nil})
			return
		}

		// make URL-safe
		id := slugify.S(resourceName.(string))

		// set in storage
		stoErr := s.Create(kind, id, payload["resource"])

		// handle error
		switch stoErr {
		case KindNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"kind not found", "", nil})
			return
		case ResourceExists:
			resp.WriteHeader(409) // Conflict
			enc.Encode(apiResponse{"resource already exists", "", nil})
			return
		}

		// report success
		resp.WriteHeader(201) // Created
		enc.Encode(apiResponse{"", id, nil})
	}

	api.Get = func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		kind := vars["kind"]
		id := vars["id"]
		enc := json.NewEncoder(resp)

		// look for resource
		resource, stoErr := s.Get(kind, id)

		// handle error
		switch stoErr {
		case KindNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"kind not found", "", nil})
			return
		case ResourceNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"resource not found", "", nil})
			return
		}

		// return artist
		err := enc.Encode(apiResponse{"", "", resource})
		if err != nil {
			log.Println(err)
		}
	}

	api.Update = func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		kind := vars["kind"]
		id := vars["id"]
		enc := json.NewEncoder(resp)
		dec := json.NewDecoder(req.Body)

		// read body and parse into interface{}
		var resource interface{}
		err := dec.Decode(&resource)

		if err != nil {
			log.Println(err)
			resp.WriteHeader(400) // Bad Request
			enc.Encode(apiResponse{"malformed json", "", nil})
			return
		}

		// update resource
		stoErr := s.Update(kind, id, resource)

		// handle error
		switch stoErr {
		case KindNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"kind not found", "", nil})
			return
		case ResourceNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"resource not found", "", nil})
			return
		}

		// 200 OK is implied
		enc.Encode(apiResponse{"", id, ""})
	}

	api.Delete = func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		kind := vars["kind"]
		id := vars["id"]
		enc := json.NewEncoder(resp)

		// delete resource
		stoErr := s.Delete(kind, id)

		// handle error
		switch stoErr {
		case KindNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"kind not found", "", nil})
			return
		case ResourceNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"resource not found", "", nil})
			return
		}

		// 200 OK is implied
		enc.Encode(apiResponse{"", "", nil})
	}

	return api
}
