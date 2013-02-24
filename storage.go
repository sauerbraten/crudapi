package crudapi

import (
	"strconv"
	"sync"
	"time"
)

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
// When implementing your own storage, make sure that at minimum the CRUD methods are thread-safe.
type Storage interface {
	// adds a new kind of resource
	AddKind(string) StorageError
	// deletes all resources of a kind, and the kind itself 
	DeleteKind(string) StorageError

	// creates a resource and stores the data in it, then returns the ID
	Create(string, map[string]interface{}) (string, StorageError)
	// retrieves a resource
	Get(string, string) (map[string]interface{}, StorageError)
	// retrieves all resources of the specified kind, returns them in a map of id → resource
	GetAll(string) (map[string]map[string]interface{}, StorageError)
	// updates a resource
	Update(string, string, map[string]interface{}) StorageError
	// deletes a resource
	Delete(string, string) StorageError
}

// MapStorage is a basic storage using maps. Thus, it is not persistent! It is meant as an example and for testing purposes.
// MapStorage is thread-safe, as any Storage implementation should be, since CRUD handlers run in parrallel as well.
type MapStorage struct {
	sync.RWMutex
	data map[string]map[string]map[string]interface{}
}

// Returns an initialized MapStorage
func NewMapStorage() MapStorage {
	return MapStorage{sync.RWMutex{}, make(map[string]map[string]map[string]interface{})}
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
	ms.data[kind] = make(map[string]map[string]interface{})
	ms.Unlock()

	return None
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

	return None
}

func (ms MapStorage) Create(kind string, resource map[string]interface{}) (id string, err StorageError) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		return
	}

	// make (pesudo-random) ID
	id = strconv.FormatInt(time.Now().Unix(), 10)

	// create nil entry for the new ID
	ms.Lock()
	ms.data[kind][id] = resource
	ms.Unlock()

	return
}

func (ms MapStorage) Get(kind, id string) (resource map[string]interface{}, err StorageError) {
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

func (ms MapStorage) GetAll(kind string) (resources map[string]map[string]interface{}, err StorageError) {
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
	resources = ms.data[kind]
	ms.RUnlock()

	return
}

func (ms MapStorage) Update(kind, id string, resource map[string]interface{}) StorageError {
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

	return None
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

	return None
}
