package main

import (
	"github.com/sauerbraten/crudapi"
	"net/http"
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
func (ms MapStorage) AddMap(kind string) {
	ms.Lock()
	ms.data[kind] = make(map[string]interface{})
	ms.Unlock()
}

// Reverts AddMap().
func (ms MapStorage) DeleteMap(kind string) {
	ms.Lock()
	delete(ms.data, kind)
	ms.Unlock()
}

func (ms MapStorage) Create(kind string, resource interface{}) (id string, resp crudapi.StorageResponse) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "kind not found"
		return
	}

	// make (pesudo-random) ID
	id = strconv.FormatInt(time.Now().Unix(), 10)

	// create nil entry for the new ID
	ms.Lock()
	ms.data[kind][id] = resource
	ms.Unlock()

	resp.StatusCode = http.StatusCreated
	return
}

func (ms MapStorage) Get(kind, id string) (resource interface{}, resp crudapi.StorageResponse) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "kind not found"
		return
	}

	// make sure a resource with this ID exists
	ms.RLock()
	resource, ok = ms.data[kind][id]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "resource not found"
		return
	}

	resp.StatusCode = http.StatusOK
	return
}

func (ms MapStorage) GetAll(kind string) (resources []interface{}, resp crudapi.StorageResponse) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "kind not found"
		return
	}

	// collect all values in the kind's map in a slice
	ms.RLock()
	for _, resource := range ms.data[kind] {
		resources = append(resources, resource)
	}
	ms.RUnlock()

	resp.StatusCode = http.StatusOK
	return
}

func (ms MapStorage) Update(kind, id string, resource interface{}) (resp crudapi.StorageResponse) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "kind not found"
		return
	}

	// make sure the resource exists
	ms.RLock()
	_, ok = ms.data[kind][id]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "resource not found"
		return
	}

	// update resource
	ms.Lock()
	ms.data[kind][id] = resource
	ms.Unlock()

	resp.StatusCode = http.StatusOK
	return
}

func (ms MapStorage) Delete(kind, id string) (resp crudapi.StorageResponse) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "kind not found"
		return
	}

	// make sure the resource exists
	ms.RLock()
	_, ok = ms.data[kind][id]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "resource not found"
		return
	}

	// delete resource
	ms.Lock()
	delete(ms.data[kind], id)
	ms.Unlock()

	resp.StatusCode = http.StatusOK
	return
}

func (ms MapStorage) DeleteAll(kind string) (resp crudapi.StorageResponse) {
	// make sure kind exists
	ms.RLock()
	_, ok := ms.data[kind]
	ms.RUnlock()
	if !ok {
		resp.StatusCode = http.StatusNotFound
		resp.ErrorMessage = "kind not found"
		return
	}

	// delete resources
	ms.Lock()
	for id := range ms.data[kind] {
		delete(ms.data[kind], id)
	}
	ms.Unlock()

	resp.StatusCode = http.StatusOK
	return
}
