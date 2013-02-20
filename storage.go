package crudapi

import (
	"sync"
)

// A StorageError is returned by Storage's methods and describes what kind of error occured.
type StorageError int

const (
	_                StorageError = iota // 0 is not used and means no error
	ResourceExists                       // resource already exists
	ResourceNotFound                     // resource not found / no such resource
	KindExists                           // kind already exists
	KindNotFound                         // kind not found /no such kind
)

// Storage describes the methods required for a storage to be used with the API type.
// When implementing your own storage, make sure that at minimum the CRUD methods are thread-safe.
type Storage interface {
	AddKind(string) StorageError                     // adds a new kind of resource
	DeleteKind(string) StorageError                  // deletes all resources of a kind, and the kind itself
	Create(string, string, interface{}) StorageError // creates a resource
	Get(string, string) (interface{}, StorageError)  // retrieves a resource
	GetAll(string) ([]interface{}, StorageError)     // retrieves all resources of the specified kind
	Update(string, string, interface{}) StorageError // updates a resource
	Delete(string, string) StorageError              // deletes a resource
}

// MapStorage is a basic API storage using maps. Thus, it is not persistent! It is meant as an example and for testing purposes.
// MapStorage is thread-safe, as any Storage implementation should be, since CRUD handlers run in parrallel as well.
type MapStorage struct {
	sync.RWMutex
	data map[string]map[string]interface{}
}

// Returns an initialized MapStorage
func NewMapStorage() MapStorage {
	return MapStorage{sync.RWMutex{}, make(map[string]map[string]interface{})}
}

func (ms MapStorage) AddKind(kind string) StorageError {
	// check if kind already exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if ok {
		return KindExists
	}

	ms.Lock()
	ms.data[kind] = make(map[string]interface{})
	ms.Unlock()

	return 0
}

func (ms MapStorage) DeleteKind(kind string) StorageError {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		return KindNotFound
	}

	ms.Lock()
	delete(ms.data, kind)
	ms.Unlock()

	return 0
}

func (ms MapStorage) Create(kind, id string, resource interface{}) StorageError {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		return KindNotFound
	}

	// make sure a resource with this ID does not exist already
	ms.RLock()
	_, ok = ms.data[kind][id]
	ms.RUnlock()
	if ok {
		return ResourceExists
	}

	// create nil entry for the new ID
	ms.Lock()
	ms.data[kind][id] = resource
	ms.Unlock()

	return 0
}

func (ms MapStorage) Get(kind, id string) (resource interface{}, err StorageError) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		err = KindNotFound
		return
	}

	// make sure a resource with this ID exists
	ms.RLock()
	resource, ok = ms.data[kind][id]
	ms.RUnlock()
	if !ok {
		err = ResourceNotFound
		return
	}

	return
}

func (ms MapStorage) GetAll(kind string) (resources []interface{}, err StorageError) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		err = KindNotFound
		return
	}

	// collect all values in the kind's map in a slice
	ms.RLock()
	for _, resource := range ms.data[kind] {
		resources = append(resources, resource)
	}
	ms.RUnlock()

	return
}

func (ms MapStorage) Update(kind, id string, resource interface{}) StorageError {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		return KindNotFound
	}

	// make sure the resource exists
	ms.RLock()
	_, ok = ms.data[kind][id]
	ms.RUnlock()
	if !ok {
		return ResourceNotFound
	}

	// update resource
	ms.Lock()
	ms.data[kind][id] = resource
	ms.Unlock()

	return 0
}

func (ms MapStorage) Delete(kind, id string) StorageError {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		return KindNotFound
	}

	// make sure the resource exists
	ms.RLock()
	_, ok = ms.data[kind][id]
	ms.RUnlock()
	if !ok {
		return ResourceNotFound
	}

	// update resource
	ms.Lock()
	delete(ms.data[kind], id)
	ms.Unlock()

	return 0
}
