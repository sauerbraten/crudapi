package crudapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

	// Create is meant to handle POST requests. It returns '400 Bad Request', '404 Not Found' or '201 Created'.
	create func(resp http.ResponseWriter, req *http.Request)

	// Get is meant to retrieve resources (HTTP GET). It returns '404 Not Found' or '200 OK'.
	get func(resp http.ResponseWriter, req *http.Request)

	// GetAll returns all resources of a kind.
	getAll func(resp http.ResponseWriter, req *http.Request)

	// Update is meant to handle PUTs. It returns '400 Bad Request', '404 Not Found' or '200 OK'.
	update func(resp http.ResponseWriter, req *http.Request)

	// Delete is for handling DELETE requests on single resources. Possible HTTP status codes are '404 Not Found' and '200 OK'.
	delete func(resp http.ResponseWriter, req *http.Request)

	// DeleteAll is for handling DELETE on whole kinds. Possible HTTP status codes are '404 Not Found' and '200 OK'.
	deleteAll func(resp http.ResponseWriter, req *http.Request)

	// The generated CRUD routes. Pass this to http.ListenAndServe().
	Router *mux.Router
}

// Returns an API relying on the given Storage. The path prefix must start with '/' and be 2 or more characters long. If those criteria are not met, "/api" is used as a fallback instead. Trailing slashes are stripped, like so: "/api///" → "/api".
func NewAPI(pathPrefix string, s Storage) *API {
	// validate path prefix
	if len(pathPrefix) < 2 || pathPrefix[0] != '/' {
		pathPrefix = "/api"
	}

	// strip trailing slashes
	for pathPrefix[len(pathPrefix)-1] == '/' {
		pathPrefix = pathPrefix[:len(pathPrefix)-1]
	}

	api := &API{}

	api.s = s

	api.create = func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		kind := vars["kind"]
		enc := json.NewEncoder(resp)
		dec := json.NewDecoder(req.Body)

		// read body and parse into interface{}
		var resource map[string]interface{}
		err := dec.Decode(&resource)

		if err != nil {
			log.Println(err)
			resp.WriteHeader(400) // Bad Request
			enc.Encode(apiResponse{"malformed json", "", nil})
			return
		}

		// set in storage
		id, stoErr := s.Create(kind, resource)

		// handle error
		switch stoErr {
		case KindNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"kind not found", "", nil})
			return
		case InternalError:
			resp.WriteHeader(500) // Internal Server Error
			enc.Encode(apiResponse{"storage failure", "", nil})
			return
		}

		// report success
		resp.WriteHeader(201) // Created
		enc.Encode(apiResponse{"", id, nil})
	}

	api.get = func(resp http.ResponseWriter, req *http.Request) {
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
		case InternalError:
			resp.WriteHeader(500) // Internal Server Error
			enc.Encode(apiResponse{"storage failure", "", nil})
			return
		}

		// return resource
		err := enc.Encode(apiResponse{"", "", resource})
		if err != nil {
			log.Println(err)
		}
	}

	api.getAll = func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		kind := vars["kind"]
		enc := json.NewEncoder(resp)

		// look for resources
		resources, stoErr := s.GetAll(kind)

		// handle error
		switch stoErr {
		case KindNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"kind not found", "", nil})
			return
		case InternalError:
			resp.WriteHeader(500) // Internal Server Error
			enc.Encode(apiResponse{"storage failure", "", nil})
			return
		}

		// return resources
		err := enc.Encode(apiResponse{"", "", resources})
		if err != nil {
			log.Println(err)
		}
	}

	api.update = func(resp http.ResponseWriter, req *http.Request) {
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
		case InternalError:
			resp.WriteHeader(500) // Internal Server Error
			enc.Encode(apiResponse{"storage failure", "", nil})
			return
		}

		// 200 OK is implied
		enc.Encode(apiResponse{"", id, nil})
	}

	api.delete = func(resp http.ResponseWriter, req *http.Request) {
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
		case InternalError:
			resp.WriteHeader(500) // Internal Server Error
			enc.Encode(apiResponse{"storage failure", "", nil})
			return
		}

		// 200 OK is implied
		enc.Encode(apiResponse{"", "", nil})
	}

	api.deleteAll = func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		kind := vars["kind"]
		enc := json.NewEncoder(resp)

		// look for resources
		stoErr := s.DeleteAll(kind)

		// handle error
		switch stoErr {
		case KindNotFound:
			resp.WriteHeader(404) // Not Found
			enc.Encode(apiResponse{"kind not found", "", nil})
			return
		case InternalError:
			resp.WriteHeader(500) // Internal Server Error
			enc.Encode(apiResponse{"storage failure", "", nil})
			return
		}

		// 200 OK is implied
		enc.Encode(apiResponse{"", "", nil})
	}

	api.Router = mux.NewRouter()
	api.Router.StrictSlash(true)

	r := api.Router.PathPrefix(pathPrefix).Subrouter()

	// set up CRUD routes
	r.HandleFunc("/{kind}", api.create).Methods("POST")
	r.HandleFunc("/{kind}", api.getAll).Methods("GET")
	r.HandleFunc("/{kind}/{id}", api.get).Methods("GET")
	r.HandleFunc("/{kind}/{id}", api.update).Methods("PUT")
	r.HandleFunc("/{kind}", api.deleteAll).Methods("DELETE")
	r.HandleFunc("/{kind}/{id}", api.delete).Methods("DELETE")

	return api
}
