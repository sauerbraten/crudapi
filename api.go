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
	"net/url"
)

type apiResponse struct {
	ErrorMessage string      `json:"error,omitempty"`
	Id           string      `json:"id,omitempty"`
	Result       interface{} `json:"result,omitempty"`
}

var s Storage
var g Guard

// Adds CRUD and OPTIONS routes to the router, which rely on the given Storage. If guard is nil, all requests are allowed by default.
func MountAPI(router *mux.Router, storage Storage, guard Guard) {
	s = storage
	if s == nil {
		panic("storage is nil")
	}

	g = guard
	if g == nil {
		g = defaultGuard{}
	}

	// Create
	router.HandleFunc("/{kind}", create).Methods("POST")

	// Read
	router.HandleFunc("/{kind}", getAll).Methods("GET")
	router.HandleFunc("/{kind}/{id}", get).Methods("GET")

	// Update
	router.HandleFunc("/{kind}/{id}", update).Methods("PUT")

	// Delete
	router.HandleFunc("/{kind}", deleteAll).Methods("DELETE")
	router.HandleFunc("/{kind}/{id}", del).Methods("DELETE")

	// OPTIONS routes for API discovery
	router.HandleFunc("/{kind}", optionsKind).Methods("OPTIONS")
	router.HandleFunc("/{kind}/{id}", optionsResource).Methods("OPTIONS")
}

// initializes variables needed to handle every request and authenticates and authorizes the request
func initHandling(req *http.Request, resp http.ResponseWriter, action Action) (authenticatedAndAuthorized bool, kind string, enc *json.Encoder) {
	vars := mux.Vars(req)
	kind = vars["kind"]
	params := req.URL.Query()
	enc = json.NewEncoder(resp)

	authenticatedAndAuthorized = authenticateAndAuthorize(action, vars["kind"], params, resp, enc)

	return
}

// authenticate request and authorize action
func authenticateAndAuthorize(action Action, kind string, params url.Values, resp http.ResponseWriter, enc *json.Encoder) (ok bool) {
	authenticated, client, errorMessage := g.Authenticate(params)
	if !authenticated {
		log.Println("unauthenticated request:\n\tURL parameters:", params, "\n\terror message:", errorMessage)

		resp.WriteHeader(http.StatusUnauthorized)
		err := enc.Encode(apiResponse{errorMessage, "", nil})
		if err != nil {
			log.Println(err)
		}

		return
	}

	authorized, errorMessage := g.Authorize(client, action, kind)
	if !authorized {
		log.Println("unauthorized request:\n\tclient:", client, "\n\taction:", action, "kind:", kind, "\n\terror message", errorMessage)

		resp.WriteHeader(http.StatusForbidden)
		err := enc.Encode(apiResponse{errorMessage, "", nil})
		if err != nil {
			log.Println(err)
		}

		return
	}

	ok = true
	return
}

func create(resp http.ResponseWriter, req *http.Request) {
	authenticatedAndAuthorized, kind, enc := initHandling(req, resp, ActionCreate)
	if !authenticatedAndAuthorized {
		return
	}

	// read body and parse into interface{}
	dec := json.NewDecoder(req.Body)
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
	err = enc.Encode(apiResponse{stoResp.ErrorMessage, id, nil})
	if err != nil {
		log.Println(err)
	}
}

func getAll(resp http.ResponseWriter, req *http.Request) {
	authenticatedAndAuthorized, kind, enc := initHandling(req, resp, ActionGetAll)
	if !authenticatedAndAuthorized {
		return
	}

	// look for resources
	resources, stoResp := s.GetAll(kind)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.ErrorMessage, "", resources})
	if err != nil {
		log.Println(err)
	}
}

func get(resp http.ResponseWriter, req *http.Request) {
	authenticatedAndAuthorized, kind, enc := initHandling(req, resp, ActionGet)
	if !authenticatedAndAuthorized {
		return
	}

	// look for resource
	id := mux.Vars(req)["id"]
	resource, stoResp := s.Get(kind, id)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.ErrorMessage, "", resource})
	if err != nil {
		log.Println(err)
	}
}

func update(resp http.ResponseWriter, req *http.Request) {
	authenticatedAndAuthorized, kind, enc := initHandling(req, resp, ActionUpdate)
	if !authenticatedAndAuthorized {
		return
	}

	// read body and parse into interface{}
	dec := json.NewDecoder(req.Body)
	var resource map[string]interface{}
	err := dec.Decode(&resource)

	if err != nil {
		log.Println(err)
		resp.WriteHeader(http.StatusBadRequest)
		err := enc.Encode(apiResponse{"malformed json", "", nil})
		if err != nil {
			log.Println(err)
		}

		return
	}

	// update resource
	id := mux.Vars(req)["id"]
	stoResp := s.Update(kind, id, resource)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err = enc.Encode(apiResponse{stoResp.ErrorMessage, "", nil})
	if err != nil {
		log.Println(err)
	}
}

func deleteAll(resp http.ResponseWriter, req *http.Request) {
	authenticatedAndAuthorized, kind, enc := initHandling(req, resp, ActionDeleteAll)
	if !authenticatedAndAuthorized {
		return
	}

	// look for resources
	stoResp := s.DeleteAll(kind)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.ErrorMessage, "", nil})
	if err != nil {
		log.Println(err)
	}
}

// delete() is a built-in function, thus del() is used here
func del(resp http.ResponseWriter, req *http.Request) {
	authenticatedAndAuthorized, kind, enc := initHandling(req, resp, ActionDelete)
	if !authenticatedAndAuthorized {
		return
	}

	// delete resource
	id := mux.Vars(req)["id"]
	stoResp := s.Delete(kind, id)

	// write response
	resp.WriteHeader(stoResp.StatusCode)
	err := enc.Encode(apiResponse{stoResp.ErrorMessage, "", nil})
	if err != nil {
		log.Println(err)
	}
}

func optionsKind(resp http.ResponseWriter, req *http.Request) {
	h := resp.Header()

	h.Add("Allow", "PUT")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}

func optionsResource(resp http.ResponseWriter, req *http.Request) {
	h := resp.Header()

	h.Add("Allow", "POST")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}
