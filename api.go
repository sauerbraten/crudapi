//Package crudapi implements HTTP handlers a minimalistic RESTful API offering Create, Read, Update, and Delete (→CRUD) handlers.
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
