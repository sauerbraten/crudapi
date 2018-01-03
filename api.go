/*
Package crudapi implements a RESTful JSON API exposing CRUD functionality relying on a custom storage.

See http://en.wikipedia.org/wiki/RESTful and http://en.wikipedia.org/wiki/Create,_read,_update_and_delete for more information.

An example can be found at: https://github.com/sauerbraten/crudapi#example
*/
package crudapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiResponse struct {
	ErrorMessage string      `json:"error,omitempty"`
	ID           string      `json:"id,omitempty"`
	Result       interface{} `json:"result,omitempty"`
}

type apiHandlerFunc func(Storage, http.ResponseWriter, map[string]string, *json.Encoder, *json.Decoder)

// Mount adds CRUD and OPTIONS routes to the router, which rely on the given Storage. You can provide a middleware function auth that authenticates and authorizes requests.
func Mount(router *mux.Router, storage Storage, auth func(http.HandlerFunc) http.HandlerFunc) {
	if storage == nil {
		panic("storage is nil")
	}

	if auth == nil {
		auth = func(f http.HandlerFunc) http.HandlerFunc { return f }
	}

	collectionHandlers := map[string]apiHandlerFunc{
		"GET":     getAll,
		"POST":    create,
		"DELETE":  deleteAll,
		"OPTIONS": optionsCollection,
	}

	resourceHandlers := map[string]apiHandlerFunc{
		"GET":     get,
		"PUT":     update,
		"DELETE":  del,
		"OPTIONS": optionsResource,
	}

	router.HandleFunc("/{collection}", auth(chooseAndInitialize(collectionHandlers, storage)))
	router.HandleFunc("/{collection}/{id}", auth(chooseAndInitialize(resourceHandlers, storage)))
}

func chooseAndInitialize(handlersByMethod map[string]apiHandlerFunc, storage Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		handler, ok := handlersByMethod[req.Method]
		if !ok {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(req)
		enc := json.NewEncoder(resp)
		dec := json.NewDecoder(req.Body)

		handler(storage, resp, vars, enc, dec)
	}
}

func create(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	// set in storage
	id, stoResp := storage.Create(vars["collection"], dec)

	// write response
	resp.WriteHeader(stoResp.StatusCode())
	err := enc.Encode(apiResponse{stoResp.Error(), id, nil})
	if err != nil {
		log.Println(err)
	}
}

func getAll(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	// look for resources
	resources, stoResp := storage.GetAll(vars["collection"], vars)

	// write response
	resp.WriteHeader(stoResp.StatusCode())
	err := enc.Encode(apiResponse{stoResp.Error(), "", resources})
	if err != nil {
		log.Println(err)
	}
}

func get(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	// look for resource
	resource, stoResp := storage.Get(vars["collection"], vars["id"], vars)

	// write response
	resp.WriteHeader(stoResp.StatusCode())
	err := enc.Encode(apiResponse{stoResp.Error(), "", resource})
	if err != nil {
		log.Println(err)
	}
}

func update(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	// update resource
	stoResp := storage.Update(vars["collection"], vars["id"], dec)

	// write response
	resp.WriteHeader(stoResp.StatusCode())
	err := enc.Encode(apiResponse{stoResp.Error(), "", nil})
	if err != nil {
		log.Println(err)
	}
}

func deleteAll(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	// look for resources
	stoResp := storage.DeleteAll(vars["collection"])

	// write response
	resp.WriteHeader(stoResp.StatusCode())
	err := enc.Encode(apiResponse{stoResp.Error(), "", nil})
	if err != nil {
		log.Println(err)
	}
}

// delete() is a built-in function, thus del() is used here
func del(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	// delete resource
	stoResp := storage.Delete(vars["collection"], vars["id"])

	// write response
	resp.WriteHeader(stoResp.StatusCode())
	err := enc.Encode(apiResponse{stoResp.Error(), "", nil})
	if err != nil {
		log.Println(err)
	}
}

func optionsCollection(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	h := resp.Header()

	h.Add("Allow", "PUT")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}

func optionsResource(storage Storage, resp http.ResponseWriter, vars map[string]string, enc *json.Encoder, dec *json.Decoder) {
	h := resp.Header()

	h.Add("Allow", "POST")
	h.Add("Allow", "GET")
	h.Add("Allow", "DELETE")
	h.Add("Allow", "OPTIONS")

	resp.WriteHeader(http.StatusOK)
}
