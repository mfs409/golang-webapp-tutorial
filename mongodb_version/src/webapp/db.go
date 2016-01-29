// All code for direct interaction with the database is in this file
package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
//	"encoding/json"
	"log"
)

// the database connection
var db *mgo.Database

// a user row from the database looks like this:
type User struct {
	ID       bson.ObjectId `bson:"_id"`
	State    int           `bson:"state"`
	Googleid string        `bson:"googleid"`
	Name     string        `bson:"name"`
	Email    string        `bson:"email"`
	Create   time.Time     `bson:"create"`
}

// open the database
func openDB() {
	var err error
	log.Println("opening database " + cfg.DbHost)
	m, err := mgo.Dial(cfg.DbHost)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("database open")
	db = m.DB(cfg.DbName)
}

// get a user's record, to make login/register decisions
func getUserById(googleId string) (*User, error) {
	u := User{}
	err := db.C("users").Find(bson.M{"googleid" : googleId}).Select(nil).One(&u)
	// NB: Findone returns an error on not found, so we need to
	//     disambiguate between DB errors and not-found errors
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		log.Println("Error querying users", err)
		return nil, err
	}
	// TODO: verify that we get good data sometimes
	return &u, nil
}

// insert a row into the user table
func addNewUser(id string, name string, email string, state int) error {
	u := User{
		ID:bson.NewObjectId(),
		State: state,
		Googleid: id,
		Name: name,
		Email: email,
		Create: time.Now(),
	}
	err := db.C("users").Insert(u)
	if err != nil {
		log.Println(err)
	}
	return err
	/*
TODO
	_, err := insertUser.Exec(state, id, name, email)
	if err != nil {
		log.Println(err)
		return err
	}
*/
}

// get all rows from DATA, return them as JSON
func getAllRows() string {
	/*
TODO
	// get all the data, return an empty blob on failure
	rows, err := selectRows.Query()
	if err != nil {
		return ""
	}
	defer rows.Close()

	// get column names
	columns, err := rows.Columns()
	if err != nil {
		return "error"
	}

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
	for i := 0; i < 6; i++ {
		ptrs[i] = &bytestreams[i]
	}

	// parse the rows, copy them into the array as string:string maps
	for rows.Next() {
		// get data into bytestreams as a bunch of byte streams
		err = rows.Scan(ptrs...)
		if err != nil {
			return "error"
		}

		// we're going to shuffle it into here
		rowAsMap := make(map[string]interface{})

		// for each column, create a string from the byte stream,
		// match it with its column name, and put it all in rowAsMap
		for i, bytes := range bytestreams {
			// if the row type isn't text, we can have trouble,
			// but this seems to work for ints (not sure about floats)
			var v interface{}
			b, ok := bytes.([]byte)
			if ok {
				v = string(b)
			} else {
				v = bytes
			}
			rowAsMap[columns[i]] = v
		}

		// send the parsed row into the table
		allData = append(allData, rowAsMap)
	}

	// marshall as JSON, then return it as a string
	jsonData, err := json.Marshal(allData)
	if err != nil {
		return "error"
	}
	return string(jsonData)
*/
	return ""
}

// get one row from DATA, return it as JSON
//
// NB: for better or worse, we're doing this in the same way as getRows(), so
// we'll skip commenting the redundant code
func getRow(id int) string {
	return ""
	/* TODO
	// get all the data, return an empty blob on failure
	rows, err := selectRow.Query(id)
	if err != nil {
		return "internal error"
	}
	defer rows.Close()

	// get column names
	columns, err := rows.Columns()
	if err != nil {
		return "internal error"
	}

	// parse the data... note the magic number 6 again...
	rows.Next()
	bytestreams := make([]interface{}, 6)
	ptrs := make([]interface{}, 6)
	for i := 0; i < 6; i++ {
		ptrs[i] = &bytestreams[i]
	}

	// get data into bytestreams as a bunch of byte streams
	err = rows.Scan(ptrs...)
	if err != nil {
		return "not found"
	}
	rowAsMap := make(map[string]interface{})
	for i, bytes := range bytestreams {
		var v interface{}
		b, ok := bytes.([]byte)
		if ok {
			v = string(b)
		} else {
			v = bytes
		}
		rowAsMap[columns[i]] = v
	}

	// marshall as JSON, then return it as a string
	jsonData, err := json.Marshal(rowAsMap)
	if err != nil {
		return "internal error"
	}
	return string(jsonData)
*/
}

// Update a row in DATA
func updateDataRow(id int, data map[string]interface{}) bool {
	return false;
	/* TODO
	// special case for trickfloat
	var special *string = nil
	tf := data["trickfloat"].(string)
	if tf != "" {
		special = &tf
	}
	// do the update
	_, err := updateRow.Exec(data["smallnote"], data["bignote"], data["favint"], data["favfloat"], special, id)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
*/
}

// Delete a row from DATA
func deleteDataRow(id int) bool {
	return false
	/*
	_, err := deleteRow.Exec(id)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
*/
}

// Insert a row into DATA
func insertDataRow(data map[string]interface{}) bool {
	return false
	/*
	// special case for trickfloat
	var special *string = nil
	tf := data["trickfloat"].(string)
	if tf != "" {
		special = &tf
	}
	// do the insert
	_, err := insertRow.Exec(data["smallnote"], data["bignote"], data["favint"], data["favfloat"], special)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
*/
}
