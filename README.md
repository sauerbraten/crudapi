# CRUD API

An implementation of a minimalist RESTful API offering Create, Read, Update, and Delete (CRUD) handlers.

For more information, check out the wikipedia aticles for [CRUD](http://en.wikipedia.org/wiki/Create,_read,_update_and_delete) and [RESTful](http://en.wikipedia.org/wiki/RESTful).

## Usage

Get the package:

	$ go get github.com/sauerbraten/crudapi

Import the package:

	import (
		"github.com/sauerbraten/crudapi"
	)

## Example

Put this code into a `main.go` file:

	package main

	import (
		"github.com/gorilla/mux"
		"github.com/sauerbraten/crudapi"
		"log"
		"net/http"
	)

	func main() {
		// storage
		s := crudapi.NewMapStorage()
		s.AddKind("artist")
		s.AddKind("album")

		api := crudapi.NewAPI(s)

		// routes
		r := mux.NewRouter()
		r.StrictSlash(true)

		post := r.Methods("POST").Subrouter()
		get := r.Methods("GET").Subrouter()
		put := r.Methods("PUT").Subrouter()
		del := r.Methods("DELETE").Subrouter()

		// crud
		post.HandleFunc("/{kind}", api.Create)
		get.HandleFunc("/{kind}/{id}", api.Get)
		put.HandleFunc("/{kind}/{id}", api.Update)
		del.HandleFunc("/{kind}/{id}", api.Delete)

		// start listening
		log.Println("server listening on localhost:8080")
		http.ListenAndServe(":8080", r)
	}

When the server is running, try the following commands:

Create *Gorillaz* as *artist*:

	curl -i -X POST -d '{"id":"Gorillaz","resource":{"name":"Gorillaz","albums":["the-fall"]}}' http://localhost:8080/artist

Output:

	HTTP/1.1 201 Created
	Date: Wed, 20 Feb 2013 13:20:13 GMT
	Transfer-Encoding: chunked
	Content-Type: text/plain; charset=utf-8

	{"id":"gorillaz"}

When POSTing a resource, you *have* to post a JSON object with `"id"` and `"resources"` fields. The `"id"` value is slugified by the server to be URL-safe, so whitespace and fancy characters aren't an issue. You are probably fine just using a unique field of your actual resource.

Create *Plastic Beach* as *album*:

	curl -i -X POST -d '{"id":"Plastic Beach","resource":{"title":"Plastic Beach","by":"gorillaz","songs":["on-melancholy-hill","stylo"]}}' http://localhost:8080/album

Output:

	HTTP/1.1 201 Created
	Date: Wed, 20 Feb 2013 13:21:14 GMT
	Transfer-Encoding: chunked
	Content-Type: text/plain; charset=utf-8

	{"id":"plastic-beach"}

Retrieve the *Gorillaz* artist object:

	curl -i -X GET http://localhost:8080/artist/gorillaz

Output:

	HTTP/1.1 200 OK
	Date: Wed, 20 Feb 2013 13:25:59 GMT
	Transfer-Encoding: chunked
	Content-Type: text/plain; charset=utf-8

	{"resource":{"name":"Gorillaz","albums":["the-fall"]}}

Update the *Gorillaz* object and add the *Plastic Beach* album:

	curl -i -X PUT -d '{"name":"Gorillaz","albums":["plastic-beach","the-fall"]}' http://localhost:8080/artist/gorillaz

Output:

	HTTP/1.1 200 OK
	Date: Wed, 20 Feb 2013 13:27:35 GMT
	Transfer-Encoding: chunked
	Content-Type: text/plain; charset=utf-8

	{"id":"gorillaz"}

Again, retrieve the *Gorillaz* artist object:

	curl -i -X GET http://localhost:8080/artist/gorillaz

Output:

	HTTP/1.1 200 OK
	Date: Wed, 20 Feb 2013 13:28:32 GMT
	Transfer-Encoding: chunked
	Content-Type: text/plain; charset=utf-8

	{"resource":{"albums":["plastic-beach","the-fall"],"name":"Gorillaz"}}


Note the **returned HTTP codes**:

- `201 Created` when POSTing,
- `200 OK` when GETting and PUTting.

There are also

- `404 Not Found` if either the kind of data you are posting (for example `artist` and `album` in the URLs) is unkown or there is no resource with the specified id ('gorillaz' in the GET request). In that case a JSON object containing an `"error"` field is returned, i.e.: `{"error":"resource not found"}` or `{"error":"kind not found"}`.
- `400 Bad Request` is returned when either the POSTed or PUTted JSON data is malformed and cannot be parsed or when you are POSTing/PUTting without an `"id"` field in the top-level JSON object.
- `409 Conflict` and `{"error":"resource already exists"}` as response means, well, that you POSTed a resource with an `"id"` that is already in use.

Server responses are always a JSON object, containing one or more of the following fields:

- `"error"`: specifies the error that occured, if any
- `"id"`: the ID of the newly created or updated resource
- `"resource"`: the requested resource (used when GETting resources)


## Documentation

Full package documentation on [GoDoc](http://godoc.org/github.com/sauerbraten/crudapi).

## License

Copyright (c) 2013 Alexander Willing. All rights reserved.

- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.