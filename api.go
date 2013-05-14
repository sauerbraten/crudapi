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

	// set up CRUD routes
	router.HandleFunc("/{kind}", create).Methods("POST")
	router.HandleFunc("/{kind}", getAll).Methods("GET")
	router.HandleFunc("/{kind}/{id}", get).Methods("GET")
	router.HandleFunc("/{kind}/{id}", update).Methods("PUT")
	router.HandleFunc("/{kind}", deleteAll).Methods("DELETE")
	router.HandleFunc("/{kind}/{id}", del).Methods("DELETE")

	// set up OPTIONS routes for API discovery
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
		enc.Encode(apiResponse{"malformed json", "", nil})
		return
	}

	// set in storage
	id, stoErr := s.Create(kind, resource)

	// handle error
	switch stoErr {
	case KindNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"kind not found", "", nil})
		return
	case InternalError:
		resp.WriteHeader(http.StatusInternalServerError)
		enc.Encode(apiResponse{"storage failure", "", nil})
		return
	}

	// report success
	resp.WriteHeader(http.StatusCreated)
	enc.Encode(apiResponse{"", id, nil})
}

func get(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	id := vars["id"]
	enc := json.NewEncoder(resp)

	// look for resource
	resource, stoErr := s.Get(kind, id)

	// handle error
	switch stoErr {
	case KindNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"kind not found", "", nil})
		return
	case ResourceNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"resource not found", "", nil})
		return
	case InternalError:
		resp.WriteHeader(http.StatusInternalServerError)
		enc.Encode(apiResponse{"storage failure", "", nil})
		return
	}

	// return resource
	err := enc.Encode(apiResponse{"", "", resource})
	if err != nil {
		log.Println(err)
	}
}

func getAll(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	enc := json.NewEncoder(resp)

	// look for resources
	resources, stoErr := s.GetAll(kind)

	// handle error
	switch stoErr {
	case KindNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"kind not found", "", nil})
		return
	case InternalError:
		resp.WriteHeader(http.StatusInternalServerError)
		enc.Encode(apiResponse{"storage failure", "", nil})
		return
	}

	// return resources
	err := enc.Encode(apiResponse{"", "", resources})
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
		enc.Encode(apiResponse{"malformed json", "", nil})
		return
	}

	// update resource
	stoErr := s.Update(kind, id, resource)

	// handle error
	switch stoErr {
	case KindNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"kind not found", "", nil})
		return
	case ResourceNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"resource not found", "", nil})
		return
	case InternalError:
		resp.WriteHeader(http.StatusInternalServerError)
		enc.Encode(apiResponse{"storage failure", "", nil})
		return
	}

	// 200 OK is implied
	enc.Encode(apiResponse{"", "", nil})
}

// delete() is a built-in function for maps, thus the shorthand name 'del' for this handler function
func del(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	id := vars["id"]
	enc := json.NewEncoder(resp)

	// delete resource
	stoErr := s.Delete(kind, id)

	// handle error
	switch stoErr {
	case KindNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"kind not found", "", nil})
		return
	case ResourceNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"resource not found", "", nil})
		return
	case InternalError:
		resp.WriteHeader(http.StatusInternalServerError)
		enc.Encode(apiResponse{"storage failure", "", nil})
		return
	}

	// 200 OK is implied
	enc.Encode(apiResponse{"", "", nil})
}

func deleteAll(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	kind := vars["kind"]
	enc := json.NewEncoder(resp)

	// look for resources
	stoErr := s.DeleteAll(kind)

	// handle error
	switch stoErr {
	case KindNotFound:
		resp.WriteHeader(http.StatusNotFound)
		enc.Encode(apiResponse{"kind not found", "", nil})
		return
	case InternalError:
		resp.WriteHeader(http.StatusInternalServerError)
		enc.Encode(apiResponse{"storage failure", "", nil})
		return
	}

	// 200 OK is implied
	enc.Encode(apiResponse{"", "", nil})
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
