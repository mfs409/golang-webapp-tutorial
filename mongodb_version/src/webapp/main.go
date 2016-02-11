// A demo web application to show how to use OAuth 2.0 (Google+ Provider) and
// MongoDB from Go.
package main

import (
	"flag"
	"log"
	"net/http"
)

// The main function just configures resources and starts listening for new
// web requests
func main() {
	// parse command line options
	configPath := flag.String("configfile", "config.json", "Path to the configuration (JSON) file")
	flag.Parse()

	// load the JSON config file
	loadConfig(*configPath)

	// open the database
	openDB()
	defer closeDB()

	// set up templates
	buildTemplates()

	// set up the routes... it's good to have these all in one place,
	// since we need to be cautious about orders when there is a common
	// prefix
	router := new(Router)
	// REST routes for the DATA table
	router.Register("/data/[0-9a-z]+$", "PUT", handlePutData)
	router.Register("/data/[0-9a-z]+$", "GET", handleGetDataOne)
	router.Register("/data/[0-9a-z]+$", "DELETE", handleDeleteData)
	router.Register("/data$", "POST", handlePostData)
	router.Register("/data$", "GET", handleGetAllData)
	// OAuth and login/out routes
	router.Register("/auth/google/callback$", "GET", handleGoogleCallback)
	router.Register("/register", "GET", handleGoogleRegister)
	router.Register("/logout", "GET", handleLogout)
	router.Register("/login", "GET", handleGoogleLogin)
	// Static files
	router.Register("/public/", "GET", handlePublicFile)   // NB: regexp
	router.Register("/private/", "GET", handlePrivateFile) // NB: regexp
	// The logged-in main page
	router.Register("/app", "GET", handleApp)
	// The not-logged-in main page
	router.Register("/", "GET", handleMain)

	// print a diagnostic message and start the server
	log.Println("Server running on port " + cfg.AppPort)
	http.ListenAndServe(":"+cfg.AppPort, router)
}
