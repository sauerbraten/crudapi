package main

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sauerbraten/crudapi"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type mapStorageStatusResponse struct {
	error
	statusCode int
}

func newMSSR(statusCode int, err string) *mapStorageStatusResponse {
	return &mapStorageStatusResponse{
		error:      errors.New(err),
		statusCode: statusCode,
	}
}

func (mssr *mapStorageStatusResponse) StatusCode() int {
	return mssr.statusCode
}

// MapStorage is a basic storage using maps. Thus, it is not persistent! It is meant as an example and for testing purposes.
// MapStorage is thread-safe, as any Storage implementation should be, since CRUD handlers run in parrallel as well.
type MapStorage struct {
	*sync.RWMutex
	data map[string]map[string]interface{}
}

// Returns an initialized MapStorage
func NewMapStorage() *MapStorage {
	return &MapStorage{&sync.RWMutex{}, make(map[string]map[string]interface{})}
}

// Adds a interface{} to the root level map. Equivalent to a database table.
func (ms *MapStorage) AddMap(collection string) {
	ms.Lock()
	ms.data[collection] = make(map[string]interface{})
	ms.Unlock()
}

// Reverts AddMap().
func (ms *MapStorage) DeleteMap(collection string) {
	ms.Lock()
	delete(ms.data, collection)
	ms.Unlock()
}

func (ms *MapStorage) ensureCollectionExists(collection string, resp crudapi.StorageStatusResponse) bool {
	ms.RLock()
	_, ok := ms.data[collection]
	ms.RUnlock()

	if ok {
		return true
	}

	resp = newMSSR(http.StatusNotFound, "collection not found")
	return false
}

func (ms *MapStorage) ensureResourceExists(collection, id string, resp crudapi.StorageStatusResponse) (resource interface{}, ok bool) {
	if ok = ms.ensureCollectionExists(collection, resp); !ok {
		return
	}

	ms.RLock()
	resource, ok = ms.data[collection][id]
	ms.RUnlock()

	if !ok {
		resp = newMSSR(http.StatusNotFound, "resource not found")
	}

	return
}

func (ms *MapStorage) Create(collection string, resource interface{}) (id string, resp crudapi.StorageStatusResponse) {
	// make sure collection exists
	if ok := ms.ensureCollectionExists(collection, resp); !ok {
		return
	}

	// make (pesudo-random) ID
	id = strconv.FormatInt(rand.Int63(), 10)

	// create nil entry for the new ID
	ms.Lock()
	ms.data[collection][id] = resource
	ms.Unlock()

	resp = newMSSR(http.StatusCreated, "")
	return
}

func (ms *MapStorage) Get(collection, id string) (resource interface{}, resp crudapi.StorageStatusResponse) {
	// make sure resource exists
	var ok bool
	resource, ok = ms.ensureResourceExists(collection, id, resp)
	if !ok {
		return
	}

	resp = newMSSR(http.StatusOK, "")
	return
}

func (ms *MapStorage) GetAll(collection string) (resources []interface{}, resp crudapi.StorageStatusResponse) {
	// make sure collection exists
	if ok := ms.ensureCollectionExists(collection, resp); !ok {
		return
	}

	// collect all values in the collection's map in a slice
	ms.RLock()
	for _, resource := range ms.data[collection] {
		resources = append(resources, resource)
	}
	ms.RUnlock()

	resp = newMSSR(http.StatusOK, "")
	return
}

func (ms *MapStorage) Update(collection, id string, resource interface{}) (resp crudapi.StorageStatusResponse) {
	// make sure resource exists
	if _, ok := ms.ensureResourceExists(collection, id, resp); !ok {
		return
	}

	// update resource
	ms.Lock()
	ms.data[collection][id] = resource
	ms.Unlock()

	resp = newMSSR(http.StatusOK, "")
	return
}

func (ms *MapStorage) Delete(collection, id string) (resp crudapi.StorageStatusResponse) {
	// make sure resource exists
	if _, ok := ms.ensureResourceExists(collection, id, resp); !ok {
		return
	}

	// delete resource
	ms.Lock()
	delete(ms.data[collection], id)
	ms.Unlock()

	resp = newMSSR(http.StatusOK, "")
	return
}

func (ms *MapStorage) DeleteAll(collection string) (resp crudapi.StorageStatusResponse) {
	// make sure collection exists
	if ok := ms.ensureCollectionExists(collection, resp); !ok {
		return
	}

	// delete resources
	ms.Lock()
	for id := range ms.data[collection] {
		delete(ms.data[collection], id)
	}
	ms.Unlock()

	resp = newMSSR(http.StatusOK, "")
	return
}
