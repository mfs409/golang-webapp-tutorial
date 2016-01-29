// All code for direct interaction with the database is in this file
package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"encoding/json"
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

// The type for data in the "data" table
type DataRow struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	SmallNote  string        `bson:"smallnote" json:"smallnote"`
	BigNote    string        `bson:"bignote" json:"bignote"`
	FavInt     int           `bson:"favint" json:"favint"`
	FavFloat   float64       `bson:"favfloat" json:"favfloat"`
	TrickFloat *float64      `bson:"trickfloat" json:"trickfloat"`
	Create     time.Time     `bson:"create"`
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

// close the database
// We can defer this from the main app
func closeDB() {
	log.Println("closing database")
	db.Session.Close()
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
}

// get all rows from DATA, return them as JSON
func getAllRows() string {
	var results []DataRow
	err := db.C("data").Find(nil).Sort("create").All(&results)
	if err != nil {
		log.Fatal(err)
	}
	
	// marshall as JSON, then return it as a string
	jsonData, err := json.Marshal(results)
	if err != nil {
		return "error"
	}
	return string(jsonData)
}

// get one row from DATA, return it as JSON
//
// NB: for better or worse, we're doing this in the same way as getRows(), so
// we'll skip commenting the redundant code
func getRow(id string) string {


	d := DataRow{}
	err := db.C("users").Find(bson.M{"_id" : bson.ObjectIdHex(id)}).Select(nil).One(&d)
	// NB: Findone returns an error on not found, so we need to
	//     disambiguate between DB errors and not-found errors
	if err != nil {
		if err.Error() == "not found" {
			return "not found"
		}
		log.Println("Error querying users", err)
		return "internal error"
	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		return "internal error"
	}
	return string(jsonData)
}

// Update a row in DATA
func updateDataRow(id string, data DataRow) bool {
	q := bson.M{"_id" : bson.ObjectIdHex(id)}
	fields := bson.M{"smallnote" : data.SmallNote,
		"bignote" : data.BigNote,"favint" : data.FavInt,
		"favfloat" : data.FavFloat, "trickfloat" : nil}
	if data.TrickFloat != nil {
		fields["trickfloat"] = data.TrickFloat
	}
	change := bson.M{"$set" : fields}
	err := db.C("data").Update(q, change)
	if err != nil { log.Println(err); return false }
	return true;
}

// Delete a row from DATA
func deleteDataRow(id string) bool {
	q := bson.M{"_id" : bson.ObjectIdHex(id)}
	err := db.C("data").Remove(q)
	if err != nil { log.Println(err); return false }
	return true
}

// Insert a row into DATA
func insertDataRow(data DataRow) bool {
	data.ID = bson.NewObjectId()
	err := db.C("data").Insert(data)
	if err != nil { log.Println(err); return false }
	return true
}
