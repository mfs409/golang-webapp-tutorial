// When you need to serve a folder's contents via HTTP, this is a
// satisfactory technique.  It's not super-resilient, but it is about as easy
// as the corresponding Python one-liner, without requiring us to also have
// Python installed.
//
// This is not intended for use in a production setting.  EVER.  It is useful
// when you've got some HTML5 that CORS forbids from rendering correctly via
// file://.  In such cases, you can use this to serve it via
// http://localhost:8080
package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {

	// use a command-line flag (-p) to set the port on which to serve
	port := flag.String("p", "8080",
		"The port on which to install the http server")
	// use a command-line flag (-f) to specify the root folder to serve
	folder := flag.String("f", "./",
		"The folder from which to serve requests")

	flag.Parse()

	// print our configuration
	fmt.Println("Serving " + *folder + " on port " + *port)

	// On any request, we add the folder prefix and then attempt to serve
	// the file that results.  Note that this approach will display
	// folders, too.
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, *folder+req.URL.Path)
	})
	http.ListenAndServe(":"+*port, nil)
}
