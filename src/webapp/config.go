// All hard-coded app configuration is in this file, as is all code for
// interacting with config information that is stored in a JSON file.
package main

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"os"
)

// Information regarding the Google OAuth provider... most of this can't be
// set until we load the JSON file
var oauthConf = &oauth2.Config{
	ClientID:     "",
	ClientSecret: "",
	Endpoint:     google.Endpoint,
	RedirectURL:  "",
}

// For extra security, we use this random string in OAuth calls... we pass it
// to the server, and we expect to get it back when we get the reply
//
// NB: We can't re-generate this on the fly, because we need all instances of
// the server to have the same string, or else they can't all satisfy the
// same client
const oauthStateString = "FJDKSIE7S88dhjflsid83kdlsHp7TEbpg6TwHBWdJzNwYod1i5ZTbrIF5bEoO3oP"

// This is the service to which we request identifying info for a google ID
const googleIdApi = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

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
	// first, load the JSON file and parse it into /cfg/
	f, err := os.Open(cfgFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	// second, update our OAuth stuff
	oauthConf.ClientID = cfg.ClientId
	oauthConf.ClientSecret = cfg.ClientSecret
	oauthConf.Scopes = cfg.Scopes
	oauthConf.RedirectURL = cfg.RedirectUrl
}
