# CRUD API

A router-independent implementation of a minimalist RESTful JSON API offering Create, Read, Update, and Delete (CRUD) handlers.

For more information, check out the wikipedia aticles for [CRUD](http://en.wikipedia.org/wiki/Create,_read,_update_and_delete) and [RESTful](http://en.wikipedia.org/wiki/RESTful).

* [Usage](#usage)
* [Example](#example)
* [Documentation](#documentation)
* [License](#license)

```go
package main

import (
	"log"
	"net/http"

	"gopkg.in/sauerbraten/crudapi.v2"
)

func main() {
	// storage
	storage := NewStorage()

	// create CRUD API routes
	api := crudapi.New(storage)

	// mount the API
	http.Handle("/api/", http.StripPrefix("/api", api))

	// mount a custom handler if you want
	http.HandleFunc("/", hello)

	// start listening
	log.Println("server listening on localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello there!\n"))
}
```

## Usage

Get the package:

	$ go get gopkg.in/sauerbraten/crudapi.v2

Import the package:

```go
import (
	"gopkg.in/sauerbraten/crudapi.v2"
)
```

### Storage Backend

You need to specify where you want to store data. To do so, you have to implement [`crudapi.Storage`](http://godoc.org/github.com/sauerbraten/crudapi#Storage). There is an example implementation of that interface using maps.

### Authentication & Authorization

For access control, use appropriate middleware for your router.

### Routing

`crudapi.New` returns an `http.Handler`, so you can use the package with any router you want. Using the standard library, mounting the API routes could look like this:

```go
api := crudapi.New(storage)

http.Handle("/api/", http.StripPrefix("/api", api))

http.ListenAndServe(":8080", nil)
```

Using [chi](https://github.com/go-chi/chi), mounting might be done like this:

```go
r := chi.NewRouter()
api := crudapi.New(storage)

r.Mount("/api", api)

http.ListenAndServe(":8080", r)
```

Using chi or another more powerful router than the standard library's `DefaultMux` will also be helpful for namespacing the API under a version path prefix and access control.

This will create the following CRUD routes:

* `POST /{collection}`: Creates a resource of this collection and stores the data you POSTed, then returns the ID
* `GET /{collection}`: Returns all resources of this collection
* `GET /{collection}/{id}`: Returns the resource of this collection with that ID
* `PUT /{collection}/{id}`: Updates the resource of this collection with that ID
* `DELETE /{collection}`: Deletes all resources of this collection
* `DELETE /{collection}/{id}`: Deletes the resource of this collection with that ID

### HTTP status codes

The status codes are set by your `Storage` implementation; `MapStorage` for example uses the folllowing:

* `201 Created` when creating,
* `200 OK` when getting, updating and deleting.
* `404 Not Found` if either the collection of data you are POSTing to (for example `artists` and `albums` in the URLs) is unkown or you tried to get a non-existant resource (with a wrong ID). In that case `MapStorage` also sets the error, which is then returned in the JSON response, i.e.: `{"error":"resource not found"}` or `{"error":"collection not found"}`.

There are two status codes that are returned by the API handlers:

* `400 Bad Request` is returned when either the POSTed or PUTted JSON data is malformed and cannot be parsed or when you are PUTting without an `id` in the URL.
* `405 Method Not Allowed` when the HTTP method isn't supported by the endpoint, e.g. when POSTing to a specific resource instead of a collection.

Your auth middleware is responsible for sending `401 Unauthorized` or `403 Forbidden` when appropriate.

### Response Body

Server responses are always a JSON object, containing zero or more of the following fields:

* `"error"`: specifies the error that occured, if any
* `"id"`: the ID of the newly created resource (only used when POSTing)
* `"result"`: the requested resource (`GET /collection/id`) or an array of resources (`GET /collection/`)

## Example

Change into `example/` and execute `go run *.go`. When the server is running, check out the [index page](http://localhost:8080/) and try the following commands in a terminal:

### Create

Create _Gorillaz_ as _artist_:

	$ curl -i -X POST -d '{"name":"Gorillaz","albums":[]}' http://localhost:8080/api/artists
	HTTP/1.1 201 Created
	[...]

	{"id":"7774218119985532862"}

The ID in the reply is created by your storage implementation, typically a wrapper for a database, so when you insert something you get the ID of the inserted data. The MapStorage we use here simply uses random numbers (which is definitely not recommended).

Create _Plastic Beach_ as _album_:

	$ curl -i -X POST -d '{"title":"Plastic Beach","songs":["On Melancholy Hill","Stylo"]}' http://localhost:8080/api/albums
	HTTP/1.1 201 Created
	[...]

	{"id":"5972287258414936807"}

### Read

Retrieve the _Gorillaz_ artist object:

	$ curl -i -X GET http://localhost:8080/artists/api/7774218119985532862
	HTTP/1.1 200 OK
	[...]

	{"result":{"name":"Gorillaz","albums":[]}}

### Update

Update the _Gorillaz_ object and add the _Plastic Beach_ album:

	$ curl -i -X PUT -d '{"name":"Gorillaz","albums":["5972287258414936807"]}' http://localhost:8080/api/artists/7774218119985532862
	HTTP/1.1 200 OK
	[...]

Again, retrieve the _Gorillaz_ artist object:

	$ curl -i -X GET http://localhost:8080/api/artists/7774218119985532862
	HTTP/1.1 200 OK
	[...]

	{"result":{"albums":["5972287258414936807"],"name":"Gorillaz"}}

### Delete

Delete the _Gorillaz_ object:

	$ curl -i -X DELETE http://localhost:8080/api/artists/7774218119985532862
	HTTP/1.1 200 OK
	[...]

Let's check if it's really gone:

	$ curl -i -X GET http://localhost:8080/api/artists/7774218119985532862
	HTTP/1.1 404 Not Found
	[...]

	{"error":"resource not found"}

## Documentation

[Full package documentation on GoDoc.](http://godoc.org/github.com/sauerbraten/crudapi)

## License

Copyright (c) 2013-2018 The crudapi Authors. All rights reserved.

* Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
