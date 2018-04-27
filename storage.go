package crudapi // import "gopkg.in/sauerbraten/crudapi.v2"

import (
	"encoding/json"
	"net/url"
)

// A StorageStatusResponse is returned by Storage's methods. It sets the HTTP status code of the response and describes what kind of error occurred, if any.
type StorageStatusResponse interface {
	Error() string   // the error message
	StatusCode() int // the HTTP status code that is returned to the client
}

// Storage describes the methods required for a storage to be used with the API type.
// When implementing your own storage, make sure that the methods are thread-safe.
//
// When applicable, the request body is passed in as a JSON decoder which can be used to translate
// the input into arbitrary types.
// As last parameter, each function gets the URL query parameters to allow for filtering etc.
type Storage interface {
	// creates a resource and stores the data in it, then returns the ID
	Create(collection string, body *json.Decoder, query url.Values) (string, StorageStatusResponse)

	// retrieves and returns a resource
	Get(collection, id string, query url.Values) (interface{}, StorageStatusResponse)

	// retrieves and returns all resources in the specified collection
	GetAll(collection string, query url.Values) ([]interface{}, StorageStatusResponse)

	// updates a resource
	Update(collection, id string, body *json.Decoder, query url.Values) StorageStatusResponse

	// deletes a resource
	Delete(collection, id string, query url.Values) StorageStatusResponse

	// delete all resources in a collection
	DeleteAll(collection string, query url.Values) StorageStatusResponse
}
