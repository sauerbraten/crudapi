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

You need to specify where you want to store data. You can implement [crudapi.Storage](http://godoc.org/github.com/sauerbraten/crudapi#Storage) for that purpose. There is an example implementation of that interface using maps, which we will use here:

	storage := crudapi.NewMapStorage()

Make sure your storage implementation is ready to handle the kinds of data you are going to use. For example, create the tables you'll need in you database. With MapStorage you create new maps like this:

	storage.AddMap("mytype")
	storage.AddMap("myothertype")

Make sure that these are URL-safe, since you will access them as an URL path.  
Now, create the actual API and pass it a path prefix and your storage:

	api := crudapi.NewAPI("/api", s)

This will create the following routes:

- `POST /api/{kind}` – Creates a resource of this *kind* and stores the data you POSTed, then returns the ID
- `GET /api/{kind}` – Returns all resources of this *kind*
- `GET /api/{kind}/{id}` – Returns the resource of this *kind* with that *id*
- `PUT /api/{kind}/{id}` – Updates the resource of this *kind* with that *id*
- `DELETE /api/{kind}` – Deletes all resources of this *kind*
- `DELETE /api/{kind}/{id}` – Deletes the resource of this *kind* with that *id*

Last but not least, pass `api.Router` to your http server's `ListenAndServe()`, e.g.:

	http.ListenAndServe(":8080", api.Router)

You can also define additional custom handlers, like so:

	api.Router.HandleFunc("/", index)
	api.Router.HandleFunc("/search", search)

Note: You should not define additional routes starting with the API's path prefix, since those will be interpreted by the API handlers and thus won't work for you. `api.Router` uses the [gorilla mux package](http://www.gorillatoolkit.org/pkg/mux), so you can use regular expressions and fancy stuff for your paths when using [`HandleFunc()`](http://www.gorillatoolkit.org/pkg/mux#Route.HandlerFunc); for example:

	// javascript files
	api.Router.Handle("/{fn:[a-z]+\\.js}", http.FileServer(http.Dir("js")))


## Example

Put this code into a `main.go` file:

	package main

	import (
		"github.com/sauerbraten/crudapi"
		"log"
		"net/http"
	)

	func hello(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("Hello there!"))
	}

	func main() {
		// storage
		s := crudapi.NewMapStorage()
		s.AddMap("artists")
		s.AddMap("albums")

		api := crudapi.NewAPI("/api", s)

		// custom handler
		api.Router.HandleFunc("/", hello)

		// start listening
		log.Println("server listening on localhost:8080")
		err := http.ListenAndServe(":8080", api.Router)
		if err != nil {
			log.Println(err)
		}
	}

When the server is running, check out the [index page](http://localhost:8080/) and try the following commands in a terminal:

Create *Gorillaz* as *artist*:

	curl -i -X POST -d '{"name":"Gorillaz","albums":[]}' http://localhost:8080/api/artists

Output:

	HTTP/1.1 201 Created
	[...]

	{"id":"1361703578"}

The ID in the reply is created by your storage implementation, typically a wrapper for a database, so when you insert something you get the ID of the inserted data. The MapStorage we use here simply uses the unix timestamp (which is definitely not recommended!).

Create *Plastic Beach* as *album*:

	curl -i -X POST -d '{"title":"Plastic Beach","songs":["On Melancholy Hill","Stylo"]}' http://localhost:8080/api/albums

Output:

	HTTP/1.1 201 Created
	[...]

	{"id":"1361703700"}

Retrieve the *Gorillaz* artist object:

	curl -i -X GET http://localhost:8080/api/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{"resource":{"name":"Gorillaz","albums":[]}}

Update the *Gorillaz* object and add the *Plastic Beach* album:

	curl -i -X PUT -d '{"name":"Gorillaz","albums":["1361703700"]}' http://localhost:8080/api/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{}

Again, retrieve the *Gorillaz* artist object:

	curl -i -X GET http://localhost:8080/api/artists/1361703578

Output:

	HTTP/1.1 200 OK
	[...]

	{"resource":{"albums":["1361703700"],"name":"Gorillaz"}}


Note the **returned HTTP codes**:

- `201 Created` when POSTing,
- `200 OK` when GETting, PUTting and DELETEing.

There are also

- `404 Not Found` if either the kind of data you are posting (for example `artists` and `albums` in the URLs) is unkown or there is no resource with the specified ID. In that case a JSON object containing an `error` field is returned, i.e.: `{"error":"resource not found"}` or `{"error":"kind not found"}`.
- `400 Bad Request` is returned when either the POSTed or PUTted JSON data is malformed and cannot be parsed or when you are PUTting without an `id` in the URL.

Server responses are always a JSON object, containing zero or more of the following fields:

- `"error"` – specifies the error that occured, if any
- `"id"` – the ID of the newly created resource (only used when POSTing)
- `"resource"` – the requested resource (used when GETting resources)


## Documentation

Full package documentation on [GoDoc](http://godoc.org/github.com/sauerbraten/crudapi).

## License

Copyright (c) 2013 Alexander Willing. All rights reserved.

- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.