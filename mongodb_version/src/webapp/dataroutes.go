// Routes related to REST paths for accessing the DATA table
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// a helper function to send HTTP 403 / Forbidden when the user is not logged
// in
func do403(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

// Helper routine for sending JSON back to the client a bit more cleanly
func jResp(w http.ResponseWriter, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Internal Server Error:", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(string(payload)))
}

// The GET route for all rows of DATA
func handleGetAllData(w http.ResponseWriter, r *http.Request) {
	// if authentication passes, use getAllRows to get a big JSON blob to
	// send back
	if !checkLogin(r) {
		do403(w, r)
		return
	}
	w.Write([]byte(getAllRows()))
}

// The PUT route for updating a row of DATA
func handlePutData(w http.ResponseWriter, r *http.Request) {
	// check authentication
	if !checkLogin(r) {
		do403(w, r)
		return
	}

	// get the ID from the querystring
	id := r.URL.Path[6:]

	// Get the JSON blob as raw bytes, then marshal into a DataRow
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body of PUT request")
		jResp(w, "error")
		return
	}
	var d DataRow
	err = json.Unmarshal(contents, &d)
	if err != nil {
		log.Println("Error unmarshaling JSON reply", err)
		jResp(w, "error")
		return
	}

	// send the new data to the database
	ok := updateDataRow(id, d)
	if ok {
		jResp(w, "{res: 'ok'}")
	} else {
		jResp(w, "{res: 'error'}")
	}
}

// The GET route for viewing one row of DATA
func handleGetDataOne(w http.ResponseWriter, r *http.Request) {
	// check authentication
	if !checkLogin(r) {
		do403(w, r)
		return
	}

	// get the ID from the querystring
	id := r.URL.Path[6:]

	// get a big JSON blob via getRow, send it back
	w.Write([]byte(getRow(id)))
}

// The DELETE route for removing one row of DATA
func handleDeleteData(w http.ResponseWriter, r *http.Request) {
	// authenticate, then get ID from querystring
	if !checkLogin(r) {
		do403(w, r)
		return
	}

	// get the ID from the querystring
	id :=r.URL.Path[6:]

	// delete the row
	ok := deleteDataRow(id)
	if ok {
		jResp(w, "{res: 'ok'}")
	} else {
		jResp(w, "{res: 'error'}")
	}
}

// The POST route for adding a new row of DATA
func handlePostData(w http.ResponseWriter, r *http.Request) {
	// authenticate
	if !checkLogin(r) {
		do403(w, r)
		return
	}

	// Get the JSON blob as raw bytes, then marshal into a DataRow
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body of POST request")
		jResp(w, "error")
		return
	}
	var d DataRow
	err = json.Unmarshal(contents, &d)
	if err != nil {
		log.Println("Error unmarshaling JSON reply", err)
		jResp(w, "error")
		return
	}

	// insert the data
	ok := insertDataRow(d)
	if ok {
		jResp(w, "{res: 'ok'}")
	} else {
		jResp(w, "{res: 'error'}")
	}
}
