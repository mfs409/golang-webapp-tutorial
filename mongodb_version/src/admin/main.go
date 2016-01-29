// Admin app for managing aspects of the program.  The program administers an
// app according to the information provided in a config file.
// Administrative tasks include:
//   - Create or Drop the entire Database
//   - Reset the Users or Data table
//   - Populate the Data table with data from a CSV
//   - Activate new user registrations
// There is also a simple function to show how to update this file to run
// arbitrary one-off operations on the database
package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/csv"
	"encoding/json"
	"time"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Configuration information for Google OAuth, MySQL, and Memcached.  We
// parse this from a JSON config file
//
// NB: field names must start with Capital letter for JSON parse to work
//
// NB: between the field names and JSON mnemonics, it should be easy to
//     figure out what each field does
type Config struct {
	ClientId     string   `json:"OauthGoogleClientId"`
	ClientSecret string   `json:"OauthGoogleClientSecret"`
	Scopes       []string `json:"OauthGoogleScopes"`
	RedirectUrl  string   `json:"OauthGoogleRedirectUrl"`
	DbHost       string   `json:"MongoHost"`
	DbPort       string   `json:"MongoPort"`
	DbName       string   `json:"MongoDbname"`
	McdHost      string   `json:"MemcachedHost"`
	McdPort      string   `json:"MemcachedPort"`
	AppPort      string   `json:"AppPort"`
}

// The configuration information for the app we're administering
var cfg Config

// The type for data in the "users" table
type UserEntry struct {
	ID       bson.ObjectId `bson:"_id"`
	State    int           `bson:"state"`
	Googleid string        `bson:"googleid"`
	Name     string        `bson:"name"`
	Email    string        `bson:"email"`
	Create   time.Time     `bson:"create"`
}

// The type for data in the "data" table
//
// TODO: for TrickFloat, will pointer type with null value enable empty val
//       in mgo?
type DataEntry struct {
	ID         bson.ObjectId `bson:"_id"`
	SmallNote  string        `bson:"smallnote"`
	BigNote    string        `bson:"bignote"`
	FavInt     int           `bson:"favint"`
	FavFloat   float64       `bson:"favfloat"`
	TrickFloat *float64      `bson:"trickfloat"`
	Create     time.Time     `bson:"create"`
}

// Load a JSON file that has all the config information for our app, and put
// the JSON contents into the cfg variable
func loadConfig(cfgFileName string) {
	f, err := os.Open(cfgFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
}

// Reset the database by dropping the entire database, re-creating it, and
// then creating the collections
//
// NB: since the schemas for the collections are determined by the
//     insertions, this isn't very exciting
func resetDb() {
	// Connect or fail
	db, err := mgo.Dial(cfg.DbHost)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	// We'll use monotonic mode for now...
	db.SetMode(mgo.Monotonic, true)
	
	// drop the database
	err = db.DB(cfg.DbName).DropDatabase()
	if err != nil { log.Fatal(err) }

	// create the database by putting the collections into it
	//
	// NB: we must insert and then remove a row in order for the
	//     collection to actually exist
	u := db.DB(cfg.DbName).C("users")
	id := bson.NewObjectId()
	err = u.Insert(
		&UserEntry{
			ID : id,
			State: 0,
			Googleid: "",
			Name: "",
			Email: "",
			Create: time.Now()})
	if err != nil { log.Fatal(err) }
	q := bson.M{"_id" : id }
	err = u.Remove(q)
	if err != nil { log.Fatal(err) }
	
	d := db.DB(cfg.DbName).C("data")
	id = bson.NewObjectId()
	err = d.Insert(
		&DataEntry{
			ID: id,
			SmallNote: "",
			BigNote: "",
			FavInt: 72,
			FavFloat: 2.23,
			TrickFloat: nil,
			Create: time.Now()})
	if err != nil { log.Fatal(err) }
	q = bson.M{"_id": id}
	err = d.Remove(q)
	if err != nil { log.Fatal(err) }
}

// Connect to the database, and return the corresponding database object
//
// NB: will fail if the database hasn't been created
func openDB() *mgo.Database {
	// Connect or fail
	s, err := mgo.Dial(cfg.DbHost)
	if err != nil { log.Fatal(err) }
	return s.DB(cfg.DbName)
}

// Parse a CSV so that each line becomes an array of strings, and then use
// the array of strings to push a row to the data table
//
// NB: This is hard-coded for our "data" table.
func loadCsv(csvname *string, db *mgo.Database) {
	// load the csv file
	file, err := os.Open(*csvname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// get the right collection from the database
	d := db.C("data")
	
	// parse the csv, one record at a time
	reader := csv.NewReader(file)
	reader.Comma = ','
	count := 0 // count insertions, for the sake of nice output
	for {
		// get next row... exit on EOF
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		// We have an array of strings, representing the data.  We
		// need to move it into a DataEntry struct, then we can send
		// it to mgo.  Be careful about nulls!
		fi, err := strconv.Atoi(row[2])
		if err != nil { log.Fatal(err) }
		ff, err := strconv.ParseFloat(row[3], 64)
		
		var trick *float64 = nil
		if row[4] != "" {
			tf, err := strconv.ParseFloat(row[4], 64)
			if err != nil { log.Fatal(err) }
			trick = &tf
		}
		// Create the id for this record
		id := bson.NewObjectId()
		err = d.Insert(
			&DataEntry{
				ID: id,
				SmallNote: row[0],
				BigNote: row[1],
				FavInt: fi,
				FavFloat: ff,
				TrickFloat: trick,
				Create: time.Now()})
		if err != nil { log.Fatal(err) }
		count++
	}
	log.Println("Added", count, "rows")
}

// When a user registers, the new account is not active until the
// administrator activates it.  This function lists registrations that are
// not yet activated
func listNewAccounts(db *mgo.Database) {
	// get all inactive rows from the database
	var results []UserEntry
	err := db.C("users").Find(bson.M{"state":0}).Sort("create").All(&results)
	if err != nil {
		log.Fatal(err)
	}
	// print a header
	fmt.Println("New Users:")
	fmt.Println("[id googleid name email]")
	for i, v := range results {
		fmt.Println("[",i,"]", v.ID, v.Googleid, v.Name, v.Email)
	}
}

// Since this is an administrative interface, we don't need to do anything
// too fancy for the account activation: the flags include the ID to update,
// we just use it to update the database, and we don't worry about the case
// where the account is already activated
func activateAccount(db *mgo.Database, id string) {
	q := bson.M{"_id" : bson.ObjectIdHex(id)}
	change := bson.M{"$set" : bson.M{"state" : 1}}
	err := db.C("users").Update(q, change)
	if err != nil { log.Fatal(err) }

	/*
	_, err := db.Exec("UPDATE users SET state = 1 WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
*/
}

/*
// We occasionally need to do one-off queries that can't really be predicted
// ahead of time.  When that time comes, we can edit this function,
// recompile, and then run with the "oneoff" flag to do the corresponding
// action.  For now, it's hard coded to delete userid=1 from the Users table
func doOneOff(db *sql.DB) {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
}
*/

// Main routine: Use the command line options to determine what action to
// take, then forward to the appropriate function.  Since this program is for
// quick admin tasks, all of the above functions terminate on any error.
// That means we don't need error checking in this code... any function we
// call will only return on success.
func main() {
	// parse command line options
	configPath := flag.String("configfile", "config.json", "Path to the configuration (JSON) file")
	csvName := flag.String("csvfile", "data.csv", "The csv file to parse")
	opResetDb := flag.Bool("resetdb", false, "Reset the Mongo database?")
	opCsv := flag.Bool("loadcsv", false, "Load a csv into the data table?")
	opListNewReg := flag.Bool("listnewusers", false, "List new registrations?")
	opRegister := flag.String("activatenewuser", "", "Complete pending registration for a user")
//	opOneOff := flag.Bool("oneoff", false, "Run a one-off query")
	flag.Parse()

	// load the JSON config file
	loadConfig(*configPath)

	// Reset the database?
	if *opResetDb {
		resetDb()
		return
	}
	
	// open the database
	db := openDB()

	// all other ops are handled below:
	if *opCsv {
		loadCsv(csvName, db)
	}
	if *opListNewReg {
	 	listNewAccounts(db)
	}
	if *opRegister != "" {
		activateAccount(db, *opRegister)
	}
	/*
	if *opOneOff {
		doOneOff(db)
	}
*/
}
