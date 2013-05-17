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
	router.HandleFunc("/{kind}", createOne).Methods("POST")

	// Read
	router.HandleFunc("/{kind}", readAll).Methods("GET")
	router.HandleFunc("/{kind}/{id}", readOne).Methods("GET")

	// Update
	router.HandleFunc("/{kind}/{id}", updateOne).Methods("PUT")

	// Delete
	router.HandleFunc("/{kind}", deleteAll).Methods("DELETE")
	router.HandleFunc("/{kind}/{id}", deleteOne).Methods("DELETE")

	// OPTIONS routes for API discovery
	router.HandleFunc("/{kind}", optionsAll).Methods("OPTIONS")
	router.HandleFunc("/{kind}/{id}", optionsOne).Methods("OPTIONS")
}

func createOne(resp http.ResponseWriter, req *http.Request) {
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
	err = enc.Encode(apiResponse{stoResp.Err, id, nil})
	if err != nil {
		log.Println(err)
	}
}

func readAll(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	enc := json.NewEncoder(resp)

	// look for resources
	resources, stoResp := s.GetAll(kind)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err, "", resources})
	if err != nil {
		log.Println(err)
	}
}

func readOne(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	id := vars["id"]
	enc := json.NewEncoder(resp)

	// look for resource
	resource, stoResp := s.Get(kind, id)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err, "", resource})
	if err != nil {
		log.Println(err)
	}
}

func updateOne(resp http.ResponseWriter, req *http.Request) {
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
	err = enc.Encode(apiResponse{stoResp.Err, "", nil})
	if err != nil {
		log.Println(err)
	}
}

func deleteAll(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	enc := json.NewEncoder(resp)

	// look for resources
	stoResp := s.DeleteAll(kind)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err, "", nil})
	if err != nil {
		log.Println(err)
	}
}

func deleteOne(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	id := vars["id"]
	enc := json.NewEncoder(resp)

	// delete resource
	stoResp := s.Delete(kind, id)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.Err, "", nil})
	if err != nil {
		log.Println(err)
	}
}

func optionsOne(resp http.ResponseWriter, req *http.Request) {
	h := resp.Header()

	h.Add("Allow", "POST")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}

func optionsAll(resp http.ResponseWriter, req *http.Request) {
	h := resp.Header()

	h.Add("Allow", "PUT")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}
