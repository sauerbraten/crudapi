package crudapi

import (
	"strconv"
	"sync"
	"time"
)

// MapStorage is a basic storage using maps. Thus, it is not persistent! It is meant as an example and for testing purposes.
// MapStorage is thread-safe, as any Storage implementation should be, since CRUD handlers run in parrallel as well.
type MapStorage struct {
	sync.RWMutex
	data map[string]map[string]interface{}
}

// Returns an initialized MapStorage
func NewMapStorage() MapStorage {
	return MapStorage{sync.RWMutex{}, make(map[string]map[string]interface{})}
}

// Adds a interface{} to the root level map. Equivalent to a database table.
func (ms MapStorage) AddMap(kind string) StorageError {
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

	return None
}

// Reverts AddMap().
func (ms MapStorage) DeleteMap(kind string) StorageError {
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

func (ms MapStorage) Create(kind string, resource interface{}) (id string, err StorageError) {
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

	// delete resource
	ms.Lock()
	delete(ms.data[kind], id)
	ms.Unlock()

	return None
}

func (ms MapStorage) DeleteAll(kind string) StorageError {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		return KindNotFound
	}

	// delete resources
	ms.Lock()
	for id := range ms.data[kind] {
		delete(ms.data[kind], id)
	}
	ms.Unlock()

	return None
}
