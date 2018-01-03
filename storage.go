package crudapi

import (
	"encoding/json"
)

// A StorageStatusResponse is returned by Storage's methods. It sets the HTTP status code of the response and describes what kind of error occurred, if any.
// You can satisfy the embedded error interface by embedding error in your own implementation, but make sure to initialize the embedded error with an empty error message, so calling its Error method is safe.
type StorageStatusResponse interface {
	error            // the error
	StatusCode() int // the HTTP status code that is returned to the client
}

// Storage describes the methods required for a storage to be used with the API type.
// When implementing your own storage, make sure that the methods are thread-safe.
//
// When applicable, the request body is passed in as a JSON decoder which can be used to translate
// the input into arbitrary types.
// Get and GetAll receive a map containing the route parameters (returned by a call to mux.Vars())
// to allow for filtering etc.
type Storage interface {
	// creates a resource and stores the data in it, then returns the ID
	Create(collection string, body *json.Decoder) (string, StorageStatusResponse)

	// retrieves and returns a resource
	Get(collection string, id string, vars map[string]string) (interface{}, StorageStatusResponse)

	// retrieves and returns all resources in the specified collection
	GetAll(collection string, vars map[string]string) ([]interface{}, StorageStatusResponse)

	// updates a resource
	Update(collection string, id string, body *json.Decoder) StorageStatusResponse

	// deletes a resource
	Delete(collection string, id string) StorageStatusResponse

	// delete all resources in a collection
	DeleteAll(collection string) StorageStatusResponse
}
