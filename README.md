# CRUD API

An implementation of a minimalist RESTful JSON API offering Create, Read, Update, and Delete (CRUD) handlers.

For more information, check out the wikipedia aticles for [CRUD](http://en.wikipedia.org/wiki/Create,_read,_update_and_delete) and [RESTful](http://en.wikipedia.org/wiki/RESTful).

## Usage

Get the package:

	$ go get github.com/sauerbraten/crudapi

Import the package:

	import (
		"github.com/sauerbraten/crudapi"
	)

You need to specify where you want to store data. You have to implement [`crudapi.Storage`](http://godoc.org/github.com/sauerbraten/crudapi#Storage) for that purpose. There is an example implementation of that interface using maps, which we will use here:

	storage := NewMapStorage()

Make sure your storage implementation is ready to handle the kinds of data you are going to use. For example, create the tables you'll need in you database. With MapStorage you create new maps like this:

	storage.AddMap("mytype")
	storage.AddMap("myothertype")

Make sure that these are URL-safe, since you will access them as an URL path.

You can also specify who is allowed to do what with your resources. For this, there is the [`crudapi.Guard`](http://godoc.org/github.com/sauerbraten/crudapi#Guard) interface. Implementations of this interface authenticate clients, for example using API keys, and authorize their requests (you may want to offer read-only access to you API). If you pass `nil` as guard, a default guard will be used with no restrictions of any kind. This meands everyone can do everything on your API!

For example, you can use the example guard which uses a simple map to describe valid (= allowed) actions for each kind:

	guard := MapGuard{map[string][]string{
		"artists": {crudapi.ActionCreate, crudapi.ActionGet, crudapi.ActionUpdate},
		"albums":  {crudapi.ActionCreate, crudapi.ActionGet, crudapi.ActionGetAll, crudapi.ActionUpdate},
	}}

This code allows artist resources to be created, updated, and read one-by-one, and album resources to be created, updated, read one-by-one and read all at once. The example guard does not authenticate clients, though, meaning everybody can still perform those valid actions.

Next, create a `*mux.Router` (from [gorilla/mux](http://www.gorillatoolkit.org/pkg/mux)) and mount the API:

	router := mux.NewRouter()
	crudapi.MountAPI(router, storage)

You could also use a subrouter for the API to limit it to a subdomain, and use version numbers as path prefixes:

	crudapi.MountAPI(router.Host("api.domain.com").PathPrefix("/v1").Subrouter(), storage, guard)

This will create the following CRUD routes:

- `POST /{kind}`: Creates a resource of this *kind* and stores the data you POSTed, then returns the ID
- `GET /{kind}`: Returns all resources of this *kind*
- `GET /{kind}/{id}`: Returns the resource of this *kind* with that *id*
- `PUT /{kind}/{id}`: Updates the resource of this *kind* with that *id*
- `DELETE /{kind}`: Deletes all resources of this *kind*
- `DELETE /{kind}/{id}`: Deletes the resource of this *kind* with that *id*

It also adds OPTIONS routes for easy discovery of your API:

- `OPTIONS /{kind}`: Returns `Allow: POST, GET, DELETE` in an HTTP header
- `OPTIONS /{kind}/{id]`: Returns `Allow: PUT, GET, DELETE` in an HTTP header

Last but not least, pass the `*mux.Router` to your http server's `ListenAndServe()` as usual:

	http.ListenAndServe(":8080", router)

Since the API is mounted on top of your `router`, you can also define additional custom handlers, like so:

	router.HandleFunc("/", index)
	router.HandleFunc("/search", search)

This package uses the [gorilla mux package](http://www.gorillatoolkit.org/pkg/mux), so you can use regular expressions and fancy stuff for your paths when using [`HandleFunc()`](http://www.gorillatoolkit.org/pkg/mux#Route.HandlerFunc); for example:

	// javascript files
	api.Router.Handle("/{filename:[a-z]+\\.js}", http.FileServer(http.Dir("js")))


## Example

Change into `example/` and execute `go run *.go`. When the server is running, check out the [index page](http://localhost:8080/) and try the following commands in a terminal:

Create *Gorillaz* as *artist*:

	curl -i -X POST -d '{"name":"Gorillaz","albums":[]}' http://api.localhost:8080/v1/artists

Output:

	HTTP/1.1 201 Created
	[...]

	{"id":"1361703578"}

The ID in the reply is created by your storage implementation, typically a wrapper for a database, so when you insert something you get the ID of the inserted data. The MapStorage we use here simply uses the unix timestamp (which is definitely not recommended!).

Create *Plastic Beach* as *album*:

	curl -i -X POST -d '{"title":"Plastic Beach","songs":["On Melancholy Hill","Stylo"]}' http://api.localhost:8080/v1/albums

Output:

	HTTP/1.1 201 Created
	[...]

	{"id":"1361703700"}

Retrieve the *Gorillaz* artist object:

	curl -i -X GET http://api.localhost:8080/v1/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{"result":{"name":"Gorillaz","albums":[]}}

Update the *Gorillaz* object and add the *Plastic Beach* album:

	curl -i -X PUT -d '{"name":"Gorillaz","albums":["1361703700"]}' http://api.localhost:8080/v1/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{}

Again, retrieve the *Gorillaz* artist object:

	curl -i -X GET http://api.localhost:8080/v1/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{"result":{"albums":["1361703700"],"name":"Gorillaz"}}


Note the **returned HTTP codes**. Those status codes are set by your `Storage` implementation; `MapStorage` for example uses the folllowing:

- `201 Created` when creating,
- `200 OK` when getting, updating and deleting.
- `404 Not Found` if either the kind of data you are posting (for example `artists` and `albums` in the URLs) is unkown or you tried to get a non-existant resource (with a wrong ID). In that case `MapStorage` also sets the error, which is then returned in the JSON response, i.e.: `{"error":"resource not found"}` or `{"error":"kind not found"}`.

There are a few status codes that are not set by your `Storage`, but the API handlers themselves:

- `400 Bad Request` is returned when either the POSTed or PUTted JSON data is malformed and cannot be parsed or when you are PUTting without an `id` in the URL.
- `401 Unauthorized` when authentication failed. This status code is used slightly different from its original definition in the HTTP protocol, in that it does not include the `WWW-Authenticate` header field in the response, since the API does not support HTTP Basic Authentication.
- `403 Forbidden` is used when authentication was successful, but the client is not authorized to perform the specified action on the specified resource(s).

The `401 Unauthorized` status code is slightly misleading, and should actually mean `401 Unauthenticated`, since it asks for authentication, and not authorization.

Server responses are always a JSON object, containing zero or more of the following fields:

- `"error"`: specifies the error that occured, if any
- `"id"`: the ID of the newly created resource (only used when POSTing)
- `"result"`: the requested resource (`GET /kind/id`) or an array of resources (`GET /kind/`)


## Documentation

Full package documentation on [GoDoc](http://godoc.org/github.com/sauerbraten/crudapi).

## License

Copyright (c) 2013 the crudapi Authors. All rights reserved.

- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.