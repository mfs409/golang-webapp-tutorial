// Routes for handling file requests
package main

import ("net/http"; "os")

// ensure that the given path is a valid file, and not a directory
func isValidFile(path string) bool {
	file, err := os.Open(path)
	if err != nil { return false } // no file
	defer file.Close()
	stat, err := file.Stat()
	if err != nil { return false } // couldn't stat
	if stat.IsDir() { return false } // directory
	return true
}

// route for serving public static files... they are prefixed with 'public'
func handlePublicFile(w http.ResponseWriter, r *http.Request) {
	// serve only if valid file, else notfound
	path := r.URL.Path[1:]
	if isValidFile(path) {
		http.ServeFile(w, r, r.URL.Path[1:])
	} else {
		http.NotFound(w, r)
	}
}

// route for serving private static files... they are prefixed with 'private'
func handlePrivateFile(w http.ResponseWriter, r *http.Request) {
	// validate login
	if !checkLogin(r) { do403(w, r); return }
	// serve only if valid file, else notfound
	path := r.URL.Path[1:]
	if isValidFile(path) {
		http.ServeFile(w, r, r.URL.Path[1:])
	} else {
		http.NotFound(w, r)
	}
}
