/*
Package crudapi implements a RESTful JSON API exposing CRUD functionality relying on a custom storage.

See http://en.wikipedia.org/wiki/RESTful and http://en.wikipedia.org/wiki/Create,_read,_update_and_delete for more information.

An example can be found at: https://github.com/sauerbraten/crudapi#example
*/
package crudapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type apiResponse struct {
	Error  string      `json:"error,omitempty"`
	Id     string      `json:"id,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

var s Storage

// Adds CRUD and OPTIONS routes to the router, which rely on the given Storage.
func MountAPI(router *mux.Router, storage Storage) {
	s = storage

	// Create
	router.HandleFunc("/{kind}", create).Methods("POST")

	// Read
	router.HandleFunc("/{kind}", getAll).Methods("GET")
	router.HandleFunc("/{kind}/{id}", get).Methods("GET")

	// Update
	router.HandleFunc("/{kind}/{id}", update).Methods("PUT")

	// Delete
	router.HandleFunc("/{kind}", delAll).Methods("DELETE")
	router.HandleFunc("/{kind}/{id}", del).Methods("DELETE")

	// OPTIONS routes for API discovery
	router.HandleFunc("/{kind}", optionsKind).Methods("OPTIONS")
	router.HandleFunc("/{kind}/{id}", optionsId).Methods("OPTIONS")
}

func create(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	enc := json.NewEncoder(resp)
	dec := json.NewDecoder(req.Body)

	// read body and parse into interface{}
	var resource map[string]interface{}
	err := dec.Decode(&resource)

	if err != nil {
		log.Println(err)

		resp.WriteHeader(http.StatusBadRequest)
		err = enc.Encode(apiResponse{"malformed json", "", nil})
		if err != nil {
			log.Println(err)
		}

		return
	}

	// set in storage
	id, stoResp := s.Create(kind, resource)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err = enc.Encode(apiResponse{stoResp.Err.Error(), id, nil})
	if err != nil {
		log.Println(err)
	}
}

func getAll(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	enc := json.NewEncoder(resp)

	// look for resources
	resources, stoResp := s.GetAll(kind)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err.Error(), "", resources})
	if err != nil {
		log.Println(err)
	}
}

func get(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	id := vars["id"]
	enc := json.NewEncoder(resp)

	// look for resource
	resource, stoResp := s.Get(kind, id)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err.Error(), "", resource})
	if err != nil {
		log.Println(err)
	}
}

func update(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	id := vars["id"]
	enc := json.NewEncoder(resp)
	dec := json.NewDecoder(req.Body)

	// read body and parse into interface{}
	var resource map[string]interface{}
	err := dec.Decode(&resource)

	if err != nil {
		log.Println(err)
		resp.WriteHeader(http.StatusBadRequest)
		err = enc.Encode(apiResponse{"malformed json", "", nil})
		if err != nil {
			log.Println(err)
		}

		return
	}

	// update resource
	stoResp := s.Update(kind, id, resource)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err = enc.Encode(apiResponse{stoResp.Err.Error(), "", nil})
	if err != nil {
		log.Println(err)
	}
}

// delete() is a built-in function for maps, thus the shorthand name 'del' for this handler function
func del(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	id := vars["id"]
	enc := json.NewEncoder(resp)

	// delete resource
	stoResp := s.Delete(kind, id)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err.Error(), "", nil})
	if err != nil {
		log.Println(err)
	}
}

func delAll(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	enc := json.NewEncoder(resp)

	// look for resources
	stoResp := s.DeleteAll(kind)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err.Error(), "", nil})
	if err != nil {
		log.Println(err)
	}
}

func optionsKind(resp http.ResponseWriter, req *http.Request) {
	h := resp.Header()

	h.Add("Allow", "POST")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}

func optionsId(resp http.ResponseWriter, req *http.Request) {
	h := resp.Header()

	h.Add("Allow", "PUT")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}
