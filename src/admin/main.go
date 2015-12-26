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

import ("encoding/csv"; "flag"; "io"; "os"; "database/sql"; "encoding/json"
	_ "github.com/go-sql-driver/mysql"; "log"; "fmt")

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
	DbUser       string   `json:"MysqlUsername"`
	DbPass       string   `json:"MysqlPassword"`
	DbHost       string   `json:"MysqlHost"`
	DbPort       string   `json:"MysqlPort"`
	DbName       string   `json:"MysqlDbname"`
	McdHost      string   `json:"MemcachedHost"`
	McdPort      string   `json:"MemcachedPort"`
	AppPort      string   `json:"AppPort"`
}

// The configuration information for the app we're administering
var cfg Config

// Load a JSON file that has all the config information for our app, and put
// the JSON contents into the cfg variable
func loadConfig(cfgFileName string) {
	f, err := os.Open(cfgFileName)
	if err != nil { log.Fatal(err) }
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&cfg); err != nil { log.Fatal(err) }
}

// Create the database that will be used by our program.  The database name
// is read from the config file, so this code is very generic.
func createDatabase() {
	// NB: trailing '/' is necessary to indicate that we aren't
	// specifying any database
	db, err := sql.Open("mysql",
		cfg.DbUser+":"+cfg.DbPass+"@("+cfg.DbHost+":"+cfg.DbPort+")/")
	if err != nil { log.Fatal(err) }
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE `"+cfg.DbName+"`;")
	if err != nil { log.Fatal(err) }
}

// Delete the database that was being used by our program
// NB: code is almost identical to createDatabase()
func deleteDatabase() {
	db, err := sql.Open("mysql",
		cfg.DbUser+":"+cfg.DbPass+"@("+cfg.DbHost+":"+cfg.DbPort+")/")
	if err != nil { log.Fatal(err) }
	defer db.Close()
	_, err = db.Exec("DROP DATABASE `"+cfg.DbName+"`;")
	if err != nil { log.Fatal(err) }
}

// Connect to the database, and return the corresponding database object
//
// NB: will fail if the database hasn't been created
func openDB() *sql.DB {
	db, err := sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp("+cfg.DbHost+":"+cfg.DbPort+")/"+cfg.DbName)
	if err != nil { log.Fatal(err) }
	// Ping the database to be sure it's live
	err = db.Ping()
	if err != nil { log.Fatal(err) }
	return db
}

// drop and re-create the Users table.  The table name, field names, and
// field types are hard-coded into this function... for Google OAuth, you
// don't need to change this unless you want to add profile pics or other
// user-defined columns
func resetUserTable(db *sql.DB) {
	// drop the old users table
	stmt, err := db.Prepare("DROP TABLE IF EXISTS users")
	if err != nil { log.Fatal(err) }
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {log.Fatal(err) }

	// create the new users table
	//
	// NB: state is 0, 1 for pending, active
	//
	// NB: we technically don't need a prepared statement for a one-time
	//     query, but knowing the statement is legit makes debugging
	//     database interactions a tad easier
	stmt, err = db.Prepare(`CREATE TABLE users (
id       int not null auto_increment PRIMARY KEY,
state    int not null,
googleid varchar(100) not null,
name     varchar(100) not null,
email    varchar(100) not null
)`)
	if err != nil { log.Fatal(err) }
	_, err = stmt.Exec()
	if err != nil { log.Fatal(err) }
}

// drop and re-create the Data table.  The table name, field names, and field
// types are hard-coded into this function.  For any realistic app, you'll
// need more complex code, which creates multiple tables, any foreign key
// constraints, views, etc.  For our demo app, this + resetUserTable() is all
// we need to configure the data model.
func resetDataTable(db *sql.DB) {
	// drop the old table
	stmt, err := db.Prepare("DROP TABLE IF EXISTS data")
	if err != nil { log.Fatal(err) }
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil { log.Fatal(err) }

	// create the new objectives table
	stmt, err = db.Prepare(`CREATE TABLE data (
id        int not null auto_increment PRIMARY KEY,
smallnote varchar(2000) not null,
bignote   text not null,
favint    int not null,
favfloat  float not null,
trickfloat float)`)
	if err != nil { log.Fatal(err) }
	_, err = stmt.Exec()
	if err != nil { log.Fatal(err) }
}

// Parse a CSV so that each line becomes an array of strings, and then use
// the array of strings to push a row to the data table
//
// NB: It goes without saying that this is hard-coded for our Data table.
//     However, the hard-coding is quite simple: we just need to hard-code
//     the INSERT statement itself, and then the mapping of array positions
//     to VALUES in that INSERT statement.
func loadCsv(csvname *string, db *sql.DB) {
	// load the csv file
	file, err := os.Open(*csvname)
	if err != nil { log.Fatal(err) }
	defer file.Close()

	// create the prepared statement for inserting into the database
	//
	// I chose to have a single statement with 5 parameters, instead of
	// two statements to reflect the possible null-ness of the last
	// column.  That's not an objectively better approach, but for this
	// example, it results in less code.
	stmt, err := db.Prepare("INSERT INTO data(smallnote, bignote, favint, favfloat, trickfloat) VALUES(?, ?, ?, ?, ?)")
	if err != nil { log.Fatal(err) }
	
	// parse the csv, one record at a time
	reader := csv.NewReader(file)
	reader.Comma = ','
	count := 0 // count insertions, for the sake of nice output
	for {
		// get next row... exit on EOF
		row, err := reader.Read()
		if err == io.EOF { break } else if err != nil { log.Fatal(err) }
		// We have an array of strings, representing the data.  We
		// can dump it into the database by matching array indices
		// with parameters to the INSERT statement.  The only tricky
		// part is that our last column can be null, so we need a way
		// to handle null values, which are '""' in the CSV
		var lastcol *string = nil
		if row[4] != "" { lastcol = &row[4] }
		// NB: no need for casts... the MySQL driver will figure out
		//     the types
		_, err = stmt.Exec(row[0], row[1], row[2], row[3], lastcol)
		if err != nil { log.Fatal(err) }
		count++
	}
	log.Println("Added", count, "rows")
}

// When a user registers, the new account is not active until the
// administrator activates it.  This function lists registrations that are
// not yet activated
func listNewAccounts(db *sql.DB) {
	// get all inactive rows from the database
	rows, err := db.Query("SELECT * FROM users WHERE state = ?", 0)
	if err != nil { log.Fatal(err) }
	defer rows.Close()
	// scan into these vars, which get reused on each loop iteration
	var (id int; state int; googleid string; name string; email string)
	// print a header
	fmt.Println("New Users:")
	fmt.Println("[id googleid name email]")
	// print the rows
	for rows.Next() {
		err = rows.Scan(&id, &state, &googleid, &name, &email)
		if err != nil { log.Fatal(err) }
		fmt.Println(id, googleid, name, email)
	}
	// On error, rows.Next() returns false... but we still need to check
	// for errors
	err = rows.Err()
	if err != nil { log.Fatal(err) }
}

// Since this is an administrative interface, we don't need to do anything
// too fancy for the account activation: the flags include the ID to update,
// we just use it to update the database, and we don't worry about the case
// where the account is already activated
func activateAccount(db *sql.DB, id int) {
	_, err := db.Exec("UPDATE users SET state = 1 WHERE id = ?", id)
	if err != nil { log.Fatal(err) }
}

// We occasionally need to do one-off queries that can't really be predicted
// ahead of time.  When that time comes, we can edit this function,
// recompile, and then run with the "oneoff" flag to do the corresponding
// action.  For now, it's hard coded to delete userid=1 from the Users table
func doOneOff(db *sql.DB) {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", 1) 
	if err != nil { log.Fatal(err) }
}

// Main routine: Use the command line options to determine what action to
// take, then forward to the appropriate function.  Since this program is for
// quick admin tasks, all of the above functions terminate on any error.
// That means we don't need error checking in this code... any function we
// call will only return on success.
func main() {
	// parse command line options
	configPath := flag.String("configfile", "config.json", "Path to the configuration (JSON) file")
	csvName := flag.String("csvfile", "data.csv", "The csv file to parse")
	opCreateDb := flag.Bool("createschema", false, "Create the MySQL database into which we will create tables?")
	opDeleteDb := flag.Bool("deleteschema", false, "Delete the entire MySQL database?")
	opResetUserTable := flag.Bool("resetuserstable", false, "Delete and re-create the users table?")
	opResetDataTable := flag.Bool("resetdatatable", false, "Delete and re-create the data table?")
	opCsv   := flag.Bool("loadcsv", false, "Load a csv into the data table?")
	opListNewReg := flag.Bool("listnewusers", false, "List new registrations?")
	opRegister := flag.Int("activatenewuser", -1, "Complete pending registration for a user")
	opOneOff := flag.Bool("oneoff", false, "Run a one-off query")
	flag.Parse()

	// load the JSON config file
	loadConfig(*configPath)

	// if we are creating or deleting a database, do the op and return
	if *opCreateDb { createDatabase(); return }
	if *opDeleteDb { deleteDatabase(); return }

	// open the database
	db := openDB()
	defer db.Close()

	// all other ops are handled below:
	if *opResetUserTable { resetUserTable(db) }
	if *opResetDataTable { resetDataTable(db) }
	if *opCsv { loadCsv(csvName, db) }
	if *opListNewReg { listNewAccounts(db) }
	if *opRegister != -1 { activateAccount(db, * opRegister) }
	if *opOneOff { doOneOff(db) }
}
