// Package crudapi implements a RESTful JSON API exposing CRUD functionality relying on a custom storage.
//
// See http://en.wikipedia.org/wiki/RESTful and http://en.wikipedia.org/wiki/Create,_read,_update_and_delete for more information.
//
// An example can be found at: https://github.com/sauerbraten/crudapi#example
package crudapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

// New returns a handler mapping paths to the methods provided storage implementation.
// If storage is nil, New panics.
func New(storage Storage) http.Handler {
	if storage == nil {
		panic(errors.New("crudapi: storage is nil"))
	}

	h := chi.NewMux()

	routes := map[string]map[string]handler{
		"/{collection}": {
			http.MethodGet:    getAll,
			http.MethodPost:   create,
			http.MethodDelete: deleteAll,
		},
		"/{collection}/{id}": {
			http.MethodGet:    get,
			http.MethodPut:    update,
			http.MethodDelete: del,
		},
	}

	for pattern, handlers := range routes {
		for method, handler := range handlers {
			h.Method(method, pattern, handler.withStorage(storage))
		}
	}

	return h
}

type handler func(s Storage, collection, id string, query url.Values, body *json.Decoder) (statusCode int, resp response)

func (h handler) withStorage(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collection := chi.URLParam(r, "collection")
		id := chi.URLParam(r, "id")
		query := r.URL.Query()
		body := json.NewDecoder(r.Body)
		defer r.Body.Close()

		statusCode, resp := h(s, collection, id, query, body)

		w.WriteHeader(statusCode)

		if !resp.isEmpty() {
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func create(storage Storage, collection, _ string, query url.Values, body *json.Decoder) (int, response) {
	id, resp := storage.Create(collection, body, query)
	return resp.StatusCode(), response{resp.Error(), id, nil}
}

func getAll(storage Storage, collection, _ string, query url.Values, _ *json.Decoder) (int, response) {
	resources, resp := storage.GetAll(collection, query)
	return resp.StatusCode(), response{resp.Error(), "", resources}
}

func get(storage Storage, collection, id string, query url.Values, _ *json.Decoder) (int, response) {
	resource, resp := storage.Get(collection, id, query)
	return resp.StatusCode(), response{resp.Error(), "", resource}
}

func update(storage Storage, collection, id string, query url.Values, body *json.Decoder) (int, response) {
	resp := storage.Update(collection, id, body, query)
	return resp.StatusCode(), response{resp.Error(), "", nil}
}

func deleteAll(storage Storage, collection, _ string, query url.Values, _ *json.Decoder) (int, response) {
	resp := storage.DeleteAll(collection, query)
	return resp.StatusCode(), response{resp.Error(), "", nil}
}

// delete() is a built-in function, thus del() is used here
func del(storage Storage, collection, id string, query url.Values, _ *json.Decoder) (int, response) {
	resp := storage.Delete(collection, id, query)
	return resp.StatusCode(), response{resp.Error(), "", nil}
}

type response struct {
	ErrorMessage string      `json:"error,omitempty"`
	ID           string      `json:"id,omitempty"`
	Result       interface{} `json:"result,omitempty"`
}

func (r *response) isEmpty() bool {
	return r.ErrorMessage == "" &&
		r.ID == "" &&
		r.Result == nil
}
