package main

import (
	"encoding/json"
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

var (
	collectionNotFound = &mapStorageStatusResponse{
		error:      errors.New("collection not found"),
		statusCode: http.StatusNotFound,
	}
	resourceNotFound = &mapStorageStatusResponse{
		error:      errors.New("resource not found"),
		statusCode: http.StatusNotFound,
	}
)

// MapStorage is a basic storage using maps. Thus, it is not persistent! It is meant as an example and for testing purposes.
// MapStorage is thread-safe, as any Storage implementation should be, since CRUD handlers run in parallel as well.
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

func (ms *MapStorage) collectionExists(collection string) bool {
	ms.RLock()
	_, ok := ms.data[collection]
	ms.RUnlock()

	return ok
}

func (ms *MapStorage) resourceExists(collection, id string) (interface{}, bool) {
	if !ms.collectionExists(collection) {
		return nil, false
	}

	ms.RLock()
	resource, ok := ms.data[collection][id]
	ms.RUnlock()

	return resource, ok
}

func (ms *MapStorage) Create(collection string, body *json.Decoder) (string, crudapi.StorageStatusResponse) {
	// make sure collection exists
	if !ms.collectionExists(collection) {
		return "", collectionNotFound
	}

	// make (pesudo-random) ID
	id := strconv.FormatInt(rand.Int63(), 10)

	// decode JSON body
	var resource map[string]interface{}
	err := body.Decode(&resource)
	if err != nil {
		return "", &mapStorageStatusResponse{
			error:      errors.New("malformed JSON: " + err.Error()),
			statusCode: http.StatusBadRequest,
		}
	}

	// insert resource
	ms.Lock()
	ms.data[collection][id] = resource
	ms.Unlock()

	return id, newMSSR(http.StatusCreated, "")
}

func (ms *MapStorage) Get(collection, id string) (interface{}, crudapi.StorageStatusResponse) {
	// make sure resource exists
	resource, ok := ms.resourceExists(collection, id)
	if !ok {
		return nil, resourceNotFound
	}

	return resource, newMSSR(http.StatusOK, "")
}

func (ms *MapStorage) GetAll(collection string) ([]interface{}, crudapi.StorageStatusResponse) {
	// make sure collection exists
	if !ms.collectionExists(collection) {
		return nil, collectionNotFound
	}

	// collect all values in the collection's map in a slice
	var resources []interface{}
	ms.RLock()
	for _, resource := range ms.data[collection] {
		resources = append(resources, resource)
	}
	ms.RUnlock()

	return resources, newMSSR(http.StatusOK, "")
}

func (ms *MapStorage) Update(collection, id string, body *json.Decoder) crudapi.StorageStatusResponse {
	// make sure resource exists
	if _, ok := ms.resourceExists(collection, id); !ok {
		return resourceNotFound
	}

	// decode JSON body
	var resource map[string]interface{}
	err := body.Decode(&resource)
	if err != nil {
		return &mapStorageStatusResponse{
			error:      errors.New("malformed JSON: " + err.Error()),
			statusCode: http.StatusBadRequest,
		}
	}

	// update resource
	ms.Lock()
	ms.data[collection][id] = resource
	ms.Unlock()

	return newMSSR(http.StatusOK, "")
}

func (ms *MapStorage) Delete(collection, id string) crudapi.StorageStatusResponse {
	// make sure resource exists
	if _, ok := ms.resourceExists(collection, id); !ok {
		return resourceNotFound
	}

	// delete resource
	ms.Lock()
	delete(ms.data[collection], id)
	ms.Unlock()

	return newMSSR(http.StatusOK, "")
}

func (ms *MapStorage) DeleteAll(collection string) crudapi.StorageStatusResponse {
	// make sure collection exists
	if !ms.collectionExists(collection) {
		return collectionNotFound
	}

	// delete resources
	ms.Lock()
	for id := range ms.data[collection] {
		delete(ms.data[collection], id)
	}
	ms.Unlock()

	return newMSSR(http.StatusOK, "")
}
