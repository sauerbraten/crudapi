# CRUD API

An implementation of a minimalist RESTful JSON API offering Create, Read, Update, and Delete (CRUD) handlers.

For more information, check out the wikipedia aticles for [CRUD](http://en.wikipedia.org/wiki/Create,_read,_update_and_delete) and [RESTful](http://en.wikipedia.org/wiki/RESTful).

- [Usage](#usage)
	- [Storage Backend](#storage-backend)
	- [Authentication & Authorization](#authentication--authorization)
	- [Routing](#routing)
- [Example](#example)
	- [Create](#create)
	- [Read](#read)
	- [Update](#update)
	- [HTTP status codes](#http-status-codes)
	- [Response Body](#response-body)
- [Documentation](#documentation)
- [License](#license)

## Usage

Get the package:

	$ go get github.com/sauerbraten/crudapi

Import the package:

	import (
		"github.com/sauerbraten/crudapi"
	)

### Storage Backend

You need to specify where you want to store data. You have to implement [`crudapi.Storage`](http://godoc.org/github.com/sauerbraten/crudapi#Storage) for that purpose. There is an example implementation of that interface using maps, which we will use here:

	storage := NewMapStorage()

Make sure your storage implementation is ready to handle the collections of data you are going to use. For example, create the tables you'll need in you database. With MapStorage you create new maps like this:

	storage.AddMap("mytype")
	storage.AddMap("myothertype")

Make sure that these are URL-safe, since you will access them as an URL path.

### Authentication & Authorization

You can control access to resources and collections by providing a middleware function (`func(http.HandlerFunc) http.HandlerFunc`). For example:

	func auth(handler http.HandlerFunc) http.HandlerFunc {
		return func(resp http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Token")
			if token == "" {
				resp.WriteHeader(http.StatusUnauthorized)
				return
			}
			if (token != superSecretAccessToken) {
				resp.WriteHeader(http.StatusForbidden)
				return
			}

			handler(resp, req)
		}
	}

### Routing

Next, create a `*mux.Router` (from [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux)) and mount the API:

	router := mux.NewRouter()
	crudapi.Mount(router, storage, auth)

You could also use a subrouter for the API to limit it to a subdomain, and use version numbers as path prefixes:

	crudapi.Mount(router.Host("api.domain.com").PathPrefix("/v1").Subrouter(), storage, auth)

This will create the following CRUD routes:

- `POST /{collection}`: Creates a resource of this *collection* and stores the data you POSTed, then returns the ID
- `GET /{collection}`: Returns all resources of this *collection*
- `GET /{collection}/{id}`: Returns the resource of this *collection* with that *id*
- `PUT /{collection}/{id}`: Updates the resource of this *collection* with that *id*
- `DELETE /{collection}`: Deletes all resources of this *collection*
- `DELETE /{collection}/{id}`: Deletes the resource of this *collection* with that *id*

It also adds OPTIONS routes for easy discovery of your API:

- `OPTIONS /{collection}`: Returns `Allow: POST, GET, DELETE` in an HTTP header
- `OPTIONS /{collection}/{id]`: Returns `Allow: PUT, GET, DELETE` in an HTTP header

Last but not least, pass the `*mux.Router` to your http server's `ListenAndServe()` as usual:

	http.ListenAndServe(":8080", router)

Since the API is mounted on top of your `router`, you can also define additional custom handlers, like so:

	router.HandleFunc("/", index)
	router.HandleFunc("/search", search)


## Example

Change into `example/` and execute `go run *.go`. When the server is running, check out the [index page](http://localhost:8080/) and try the following commands in a terminal:

### Create

Create *Gorillaz* as *artist*:

	curl -i -X POST -d '{"name":"Gorillaz","albums":[]}' http://localhost:8080/artists

Output:

	HTTP/1.1 201 Created
	[...]

	{"id":"1361703578"}

The ID in the reply is created by your storage implementation, typically a wrapper for a database, so when you insert something you get the ID of the inserted data. The MapStorage we use here simply uses random numbers (which is definitely not recommended).

Create *Plastic Beach* as *album*:

	curl -i -X POST -d '{"title":"Plastic Beach","songs":["On Melancholy Hill","Stylo"]}' http://localhost:8080/albums

Output:

	HTTP/1.1 201 Created
	[...]

	{"id":"1361703700"}

### Read

Retrieve the *Gorillaz* artist object:

	curl -i -X GET http://localhost:8080/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{"result":{"name":"Gorillaz","albums":[]}}

### Update

Update the *Gorillaz* object and add the *Plastic Beach* album:

	curl -i -X PUT -d '{"name":"Gorillaz","albums":["1361703700"]}' http://localhost:8080/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{}

Again, retrieve the *Gorillaz* artist object:

	curl -i -X GET http://localhost:8080/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{"result":{"albums":["1361703700"],"name":"Gorillaz"}}


### HTTP status codes

Note the returned HTTP codes. Those status codes are set by your `Storage` implementation; `MapStorage` for example uses the folllowing:

- `201 Created` when creating,
- `200 OK` when getting, updating and deleting.
- `404 Not Found` if either the collection of data you are POSTing to (for example `artists` and `albums` in the URLs) is unkown or you tried to get a non-existant resource (with a wrong ID). In that case `MapStorage` also sets the error, which is then returned in the JSON response, i.e.: `{"error":"resource not found"}` or `{"error":"collection not found"}`.

There are a few status codes that are not set by your `Storage`, but the API handlers themselves:

- `400 Bad Request` is returned when either the POSTed or PUTted JSON data is malformed and cannot be parsed or when you are PUTting without an `id` in the URL.
- `405 Method Not Allowed` when the HTTP method isn't supported by the endpoint, e.g. when POSTing to a specific resource instead of a collection.

Your auth middleware is responsible for sending `401 Unauthorized` or `403 Forbidden` when appropriate.


### Response Body

Server responses are always a JSON object, containing zero or more of the following fields:

- `"error"`: specifies the error that occured, if any
- `"id"`: the ID of the newly created resource (only used when POSTing)
- `"result"`: the requested resource (`GET /collection/id`) or an array of resources (`GET /collection/`)


## Documentation

[Full package documentation on GoDoc.](http://godoc.org/github.com/sauerbraten/crudapi)

## License

Copyright (c) 2013-2016 The crudapi Authors. All rights reserved.

- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.