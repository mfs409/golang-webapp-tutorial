// The purpose of this file is to put all of the interaction with MongoDB in
// one place.  This lets us change backends without having to modify other
// parts of the code, and it also makes it easier to see how the program
// interacts with data, since it's all in one place.
package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
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

// a data row from the database looks like this
//
// NB: BSON and JSON tags are provided, so that Go will auto-marshall to/from
//     BSON and JSON as needed
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
// NB: We defer() this from main()
func closeDB() {
	log.Println("closing database")
	db.Session.Close()
}

// get a user's record, to make login/register decisions
func getUserById(googleId string) (*User, error) {
	u := User{}
	err := db.C("users").Find(bson.M{"googleid": googleId}).Select(nil).One(&u)
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
func addNewUser(googleid string, name string, email string, state int) error {
	u := User{
		ID:       bson.NewObjectId(),
		State:    state,
		Googleid: googleid,
		Name:     name,
		Email:    email,
		Create:   time.Now(),
	}
	err := db.C("users").Insert(u)
	if err != nil {
		log.Println(err)
	}
	return err
}

// get all rows from the data table, return them as JSON
func getAllRows() []byte {
	// query into an array of DataRow objects
	var results []DataRow
	err := db.C("data").Find(nil).Sort("create").All(&results)
	if err != nil {
		log.Fatal(err)
	}

	// marshall as JSON, which will produce a byte stream
	jsonData, err := json.Marshal(results)
	if err != nil {
		return []byte("error")
	}
	return jsonData
}

// get one row from the data table, return it as JSON
func getRow(id bson.ObjectId) []byte {
	d := DataRow{}
	err := db.C("users").Find(bson.M{"_id": id}).Select(nil).One(&d)
	// NB: Findone returns an error on not found, so we need to
	//     disambiguate between DB errors and not-found errors
	if err != nil {
		if err.Error() == "not found" {
			return []byte("not found")
		}
		log.Println("Error querying users", err)
		return []byte("internal error")
	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		return []byte("internal error")
	}
	return jsonData
}

// Update a row in the data table
func updateDataRow(id bson.ObjectId, data DataRow) bool {
	q := bson.M{"_id": id}
	fields := bson.M{"smallnote": data.SmallNote,
		"bignote": data.BigNote, "favint": data.FavInt,
		"favfloat": data.FavFloat, "trickfloat": nil}
	if data.TrickFloat != nil {
		fields["trickfloat"] = data.TrickFloat
	}
	change := bson.M{"$set": fields}
	err := db.C("data").Update(q, change)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Delete a row from the data table
func deleteDataRow(id bson.ObjectId) bool {
	q := bson.M{"_id": id}
	err := db.C("data").Remove(q)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Insert a row into the data table
func insertDataRow(data DataRow) bool {
	data.ID = bson.NewObjectId()
	err := db.C("data").Insert(data)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
