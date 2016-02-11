// Set up an HTTP router that can manage REST requests a little more cleanly.
// In particular, we want to route as early as possible based on
// GET/PUT/POST/DELETE verbs, and we want to use regexps to match routes.
package main

import (
	"net/http"
	"regexp"
)

// each route consists of a pattern, the HTTP verb, and a function to run
type route struct {
	rxp     *regexp.Regexp // the regexp describing the route
	verb    string         // the verb
	handler http.Handler   // the handler
}

// The router's data is just an array of routes
type Router struct {
	routes []*route
}

// extend Router with a function to register a new Route
func (this *Router) Register(regex string, verb string,
	handler func(http.ResponseWriter, *http.Request)) {
	// NB: compile the regexp before saving it
	this.routes = append(this.routes, &route{regexp.MustCompile(regex), verb,
		http.HandlerFunc(handler)})
}

// Handle a request by forwarding to the appropriate route
//
// NB: http.ListenAndServe() requires this to be called ServeHTTP
func (this *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, route := range this.routes {
		if route.verb == req.Method && route.rxp.MatchString(req.URL.Path) {
			route.handler.ServeHTTP(res, req)
			return
		}
	}
	http.NotFound(res, req)
}
