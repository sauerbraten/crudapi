package crudapi

// A StorageResponse is returned by Storage's methods. It sets the HTTP status code of the response and describes what kind of error occured, if any.
type StorageResponse struct {
	statusCode int   // the HTTP status code that is returned to the client
	err        error // the error, if any
}

// Storage describes the methods required for a storage to be used with the API type.
// When implementing your own storage, make sure that the methods are thread-safe.
type Storage interface {
	// first argument is always the kind of resource (for example the database table to use)
	// a second string argument is the resource ID
	// a interface{} is a resource (for example a JSON object or a database row with map indexes ~ column names)
	//
	// creates a resource and stores the data in it, then returns the ID
	Create(string, interface{}) (string, StorageResponse)

	// retrieves a resource
	Get(string, string) (interface{}, StorageResponse)

	// retrieves all resources of the specified kind
	GetAll(string) ([]interface{}, StorageResponse)

	// updates a resource
	Update(string, string, interface{}) StorageResponse

	// deletes a resource
	Delete(string, string) StorageResponse

	// delete all resources of a kind
	DeleteAll(string) StorageResponse
}
