package crudapi

// A StorageError is returned by Storage's methods and describes what kind of error occured.
type StorageError int

const (
	None             StorageError = iota // 0 means no error
	InternalError                        // for internal errors, e.g. database connection errors
	ResourceExists                       // resource already exists
	ResourceNotFound                     // resource not found / no such resource
	KindExists                           // kind already exists
	KindNotFound                         // kind not found /no such kind
)

// Storage describes the methods required for a storage to be used with the API type.
// When implementing your own storage, make sure that the methods are thread-safe.
type Storage interface {
	// first argument is always the kind of resource (for example the database table to use)
	// a second string argument is the resource ID
	// a interface{} is a resource (for example a JSON object or a database row with map indexes ~ column names)
	//
	// creates a resource and stores the data in it, then returns the ID
	Create(string, interface{}) (string, StorageError)

	// retrieves a resource
	Get(string, string) (interface{}, StorageError)

	// retrieves all resources of the specified kind
	GetAll(string) ([]interface{}, StorageError)

	// updates a resource
	Update(string, string, interface{}) StorageError

	// deletes a resource
	Delete(string, string) StorageError

	// delete all resources of a kind
	DeleteAll(string) StorageError
}
