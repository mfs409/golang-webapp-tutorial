// All code for direct interaction with the database is in this file
package main

import ("database/sql"; _ "github.com/go-sql-driver/mysql"; "log"; "encoding/json")

// prepared statements
var selectRows *sql.Stmt  // get all data rows
var selectRow  *sql.Stmt  // get one data row by ID
var updateRow  *sql.Stmt  // update one data row
var insertRow  *sql.Stmt  // insert one row in data table
var deleteRow  *sql.Stmt  // delete one row from data table
var selectUser *sql.Stmt  // get one user
var insertUser *sql.Stmt  // add one user

// the database connection
var db *sql.DB

// a user row from the database looks like this:
type User struct {
	id int
	state int
	googleid string
	name string
	email string
}

// open the database
func openDB() {
	// assumes we've run 'CREATE SCHEMA `cfg.MysqlDbname` ;'
	db, err := sql.Open("mysql",
		cfg.DbUser+":"+cfg.DbPass+"@("+cfg.DbHost+":"+cfg.DbPort+")/"+cfg.DbName)
	if err != nil { log.Fatal(err) }
	err = db.Ping() // ensure alive...
	if err != nil { log.Fatal(err) }

	// create prepared statements for getting and creating users
	selectUser, err = db.Prepare("SELECT * FROM users WHERE googleid = ?")
	if err != nil {log.Fatal(err) }
	insertUser, err = db.Prepare("INSERT INTO users(state, googleid, name, email) VALUES (?, ?, ?, ?)")
	if err != nil { log.Fatal(err) }

	// create prepared statements for REST routes on the 'data' table
	selectRows, err = db.Prepare("SELECT * FROM data")
	if err != nil { log.Fatal(err) }
	selectRow, err = db.Prepare("SELECT * FROM data WHERE id = ?")
	if err != nil { log.Fatal(err) }
	updateRow, err = db.Prepare("UPDATE data SET smallnote = ?, bignote = ?, favint = ?, favfloat = ?, trickfloat = ? WHERE id = ?")
	if err != nil { log.Fatal(err) }
	insertRow, err = db.Prepare("INSERT INTO data(smallnote, bignote, favint, favfloat, trickfloat) VALUES(?, ?, ?, ?, ?)")
	if err != nil { log.Fatal(err) }
	deleteRow, err = db.Prepare("DELETE FROM data WHERE id = ?")
	if err != nil { log.Fatal(err) }
}

// get a user's record, to make login/register decisions
func getUserById(googleId string) (*User, error) {
	// should never fail if DB is reachable:
	rows, err := selectUser.Query(googleId)
	if err != nil {
		log.Println("query error", err)
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() { return nil, nil } // no record exists
	
	// copy row into 'user' and return it
	var user User
	err = rows.Scan(&user.id, &user.state, &user.googleid, &user.name, &user.email)
	if err != nil {
		log.Println("scan error", err)
		return nil, err
	}
	return &user, nil
}

// insert a row into the user table
func addNewUser(id string, name string, email string, state int) error {
	_, err := insertUser.Exec(state, id, name, email)
	if err != nil { log.Println(err); return err }
	return nil
}

// get all rows from DATA, return them as JSON
func getAllRows() string {
	// get all the data, return an empty blob on failure
	rows, err := selectRows.Query()
	if err != nil { return "" }
	defer rows.Close()

	// get column names
	columns, err := rows.Columns()
	if err != nil { return "error" }

	// we ultimately want to marshal an array into a JSON string, because
	// order matters
	allData := make([]map[string]interface{}, 0)

	// When we parse a row, we need a pointer for where each column goes.
	// To save some pain, we'll have an array to hold the rows, and
	// another array of pointers to those array entries.  That way, we
	// can scan to ptrs, and then use values
	//
	// NB: '6' is a magic number, representing the number of columns
	bytestreams := make([]interface{}, 6)
	ptrs := make([]interface{}, 6)
	for i := 0; i < 6; i++ { ptrs[i] = &bytestreams[i] }

	// parse the rows, copy them into the array as string:string maps
	for rows.Next() {
		// get data into bytestreams as a bunch of byte streams
		err = rows.Scan(ptrs...)
		if err != nil { return "error" }

		// we're going to shuffle it into here
		rowAsMap := make(map[string]interface{})

		// for each column, create a string from the byte stream,
		// match it with its column name, and put it all in rowAsMap
		for i, bytes := range bytestreams {
			// if the row type isn't text, we can have trouble,
			// but this seems to work for ints (not sure about floats)
			var v interface{}
			b, ok := bytes.([]byte)
			if ok { v = string(b)} else {v = bytes }
			rowAsMap[columns[i]] = v
		}
		
		// send the parsed row into the table
		allData = append(allData, rowAsMap)
	}

	// marshall as JSON, then return it as a string
	jsonData, err := json.Marshal(allData)
	if err != nil { return "error" }
	return string(jsonData)
}

// get one row from DATA, return it as JSON
//
// NB: for better or worse, we're doing this in the same way as getRows(), so
// we'll skip commenting the redundant code
func getRow(id int) string {
	// get all the data, return an empty blob on failure
	rows, err := selectRow.Query(id)
	if err != nil { return "internal error" }
	defer rows.Close()

	// get column names
	columns, err := rows.Columns()
	if err != nil { return "internal error" }

	// parse the data... note the magic number 6 again...
	rows.Next()
	bytestreams := make([]interface{}, 6)
	ptrs := make([]interface{}, 6)
	for i := 0; i < 6; i++ { ptrs[i] = &bytestreams[i] }

	// get data into bytestreams as a bunch of byte streams
	err = rows.Scan(ptrs...)
	if err != nil { return "not found" }
	rowAsMap := make(map[string]interface{})
	for i, bytes := range bytestreams {
		var v interface{}
		b, ok := bytes.([]byte)
		if ok { v = string(b)} else {v = bytes }
		rowAsMap[columns[i]] = v
	}
	
	// marshall as JSON, then return it as a string
	jsonData, err := json.Marshal(rowAsMap)
	if err != nil { return "internal error" }
	return string(jsonData)
}

// Update a row in DATA
func updateDataRow(id int, data map[string]interface{}) bool {
	// special case for trickfloat
	var special *string = nil
	tf := data["trickfloat"].(string)
	if tf != "" { special = &tf }
	// do the update
	_, err := updateRow.Exec(data["smallnote"], data["bignote"], data["favint"], data["favfloat"], special, id)
	if err != nil { log.Println(err); return false }
	return true
}

// Delete a row from DATA
func deleteDataRow(id int) bool {
	_, err := deleteRow.Exec(id)
	if err != nil { log.Println(err); return false }
	return true
}

// Insert a row into DATA
func insertDataRow(data map[string]interface{}) bool {
	// special case for trickfloat
	var special *string = nil
	tf := data["trickfloat"].(string)
	if tf != "" { special = &tf }
	// do the insert
	_, err := insertRow.Exec(data["smallnote"], data["bignote"], data["favint"], data["favfloat"], special)
	if err != nil { log.Println(err); return false }
	return true
}
