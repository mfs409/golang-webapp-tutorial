// Routes for serving content that is generated from templates
package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// The template for generating the main page (the '/' route)
var mainPage *template.Template

// The template for generating the logged in page (the '/app' route)
var appPage *template.Template

// Flash structs are used to send a message, via a cookie, back to the
// browser when we redirect home on login/logout/registration events.
type Flash struct {
	Inf     bool   // true if there's an info message
	InfText string // info message to print
	Err     bool   // true if there's an error message
	ErrText string // error message to print
}

// We call this in order to initialize all templates when we start the app
func buildTemplates() {
	var err error
	mainPage, err = template.ParseFiles("templates/main.tpl")
	if err != nil {
		log.Fatal("main template parse error", err)
	}
	appPage, err = template.ParseFiles("templates/app.tpl")
	if err != nil {
		log.Fatal("app template parse error", err)
	}
}

// The route for '/' checks for flash cookies (i == Info; e == Error) and
// uses them when generating content via the main template
func handleMain(w http.ResponseWriter, r *http.Request) {
	// Prepare to consume flash messages
	flash := Flash{false, "", false, ""}

	// if we have an 'iflash' cookie, grab its contents then erase it
	cookie, err := r.Cookie("iflash")
	if err != http.ErrNoCookie {
		flash.InfText = cookie.Value
		flash.Inf = true
		http.SetCookie(w, &http.Cookie{Name: "iflash", Value: "-1", Expires: time.Now(), Path: "/"})
	}

	// if we have an 'eflash' cookie, grab its contents then erase it
	cookie, err = r.Cookie("eflash")
	if err != http.ErrNoCookie {
		flash.ErrText = cookie.Value
		flash.Err = true
		http.SetCookie(w, &http.Cookie{Name: "eflash", Value: "-1", Expires: time.Now(), Path: "/"})
	}

	// Render the template
	mainPage.Execute(w, flash)
}

// The route for '/app' ensures the user is logged in, and then renders the
// app page via a template
func handleApp(w http.ResponseWriter, r *http.Request) {
	if !checkLogin(r) {
		do403(w, r)
		return
	}
	appPage.Execute(w, nil)
}
