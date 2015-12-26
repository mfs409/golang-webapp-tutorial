// Functions and objects that are used in order to authenticate the user
package main

import ("net/http"; "crypto/rand"; "encoding/base64"; "io"; "time"; "log";
	"golang.org/x/oauth2"; "io/ioutil"; "encoding/json"; "sync")

// create an "enum" for dealing with responses from a login request
const (
	registerOk = iota       // registration success
	registerErr = iota      // unspecified registration error
	registerExist = iota    // register a registered user
	loginOk = iota          // login success
	loginErr = iota         // unspecified login error
	loginNotReg = iota      // log in of unregistered user
	loginNotActive = iota   // log in of unactivated registered user
)

// When the user logs in, we get the Google User ID, Email, and Name.
//
// When we get the ID, it's from Google, so we can trust it.  We also put it
// in a cookie on the client browser, so that subsequent requests can prove
// their identity.
//
// The problem is that anyone who knows the ID can spoof the user.  To
// secure, we don't just save the ID, we also save a token that we randomly
// generate upon login.  If the ID and Token match, then we trust you to be
// who you say you are.
//
// Note that an attacker can still assume your identity if it can access your
// cookies, but it can't assume your identity just by knowing your ID
//
// Note, too, that multiple requests can access this map, so we need it to be
// synchronized
//
// Lastly, note that if you have multiple servers running, and a user
// migrates among servers, their login info will be lost, and they'll have to
// re-log in.  The only way to avoid that is to persist this map to some
// location that is global across nodes, and we're not going to do that in
// this simple example.
var cookieStore = struct{
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// Check if a request is being made from an authenticated context
func checkLogin(r *http.Request) (bool) {
	// grab the "id" cookie, fail if it doesn't exist
	cookie, err := r.Cookie("id")
	if err == http.ErrNoCookie { return false }

	// grab the "key" cookie, fail if it doesn't exist
	key, err := r.Cookie("key")
	if err == http.ErrNoCookie { return false }

	// make sure we've got the right stuff in the hash
	cookieStore.RLock()
	defer cookieStore.RUnlock()
	return cookieStore.m[cookie.Value] == key.Value
}

// Generate 256 bits of randomness
func sessionId() string {
    b := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, b); err != nil { return "" }
    return base64.URLEncoding.EncodeToString(b)
}

// To in order to log the user out, we need to remove the corresponding
// cookie from our local cookie store, and then erase the cookies on the
// client browser.
//
// WARNING: you must call this before any other code in the logout route, or
//	    else there is a risk that the header will already be sent.
func processLogoutRequest(w http.ResponseWriter, r *http.Request) {
	// grab the "ID" cookie, erase from map if it is found
	id, err := r.Cookie("id")
	if err != http.ErrNoCookie {
		cookieStore.Lock()
		delete(cookieStore.m, id.Value)
		cookieStore.Unlock()
		// create a log-out (info) flash
		flash := http.Cookie{Name: "iflash", Value: "Logout successful", Path: "/"}
		http.SetCookie(w, &flash)
	}
	
	// clear the cookies on the client
	clearID := http.Cookie{Name: "id", Value: "-1", Expires: time.Now(), Path: "/"}
	http.SetCookie(w, &clearID)
	clearVal := http.Cookie{Name: "key", Value: "-1", Expires: time.Now(), Path: "/"}
	http.SetCookie(w, &clearVal)
}

// This is used in the third and final step of the OAuth dance.  Google is
// sending back a URL whose QueryString encodes a "state" and a "code".  The
// "state" is something we sent to Google, that we expect to get back... it
// helps prevent spoofing.  The "code" is something we can send to Google to
// get the user's info.  From there, we save some info in a session cookie to
// keep the user logged in.
//
// NB: return values will be based on the "enum" at the top of this file
func processLoginReply(w http.ResponseWriter, r *http.Request) int {
	// extract the code and state from the querystring
	code := r.FormValue("code")
	state := r.FormValue("state")

	// choose a default error code to return, depending on login or
	// register attempt.  First character of the 'state' string is 'r'
	// for register, 'l' for login.
	errorCode := loginErr
	if (state[:1] == "r") { errorCode = registerErr }

	// validate state... it needs to match the secret we sent
	if state[1:] != oauthStateString {
		log.Println("state didn't match", oauthStateString);
		return errorCode
	}

	// convert the authorization code into a token
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("token exchange error", code);
		return errorCode
	}

	// send the token to Google, get a JSON blob back
	response, err := http.Get(googleIdApi + token.AccessToken)
	if err != nil {
		log.Println("token lookup error", response)
		return errorCode
	}

	// the JSON blob has Google ID, Name, and Email... convert to a map
	//
	// NB: we don't convert via jsonParser.Decode, because we don't know
	//     exactly what fields we'll get back
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading JSON reply")
		return errorCode
	}
	var f interface{}
	err = json.Unmarshal([]byte(string(contents)), &f)
	if err != nil {
		log.Println("Error unmarshaling JSON reply")
		return errorCode
	}
	m := f.(map[string]interface{})
	// NB: m should now hold the following interfaces:
	//   m["id"], m["email"], m["given_name"], m["name"], and m["family_name"]

	// look up the user in the database
	//
	// NB: check err first, otherwise we might mistake a 'nil' for
	//     'unregistered' when we should return registerExist
	u, err := getUserById(m["id"].(string))
	if err != nil {
		log.Println("Unspecified SQL error during user lookup")
		return errorCode
	}
	if u == nil {
		// no user... let's hope this is a registration request
		if (state[:1] != "r") {
			log.Println("Attempt to log in an unregistered user")
			return loginNotReg
		}
		// add a registration (0 == not active)
		err = addNewUser(m["id"].(string), m["name"].(string), m["email"].(string), 0)
		if err != nil {
			log.Println("error adding new user")
			return registerErr
		}
		return registerOk
	} else {
		// we have a user... let's hope this is a login request
		if (state [:1] == "r") {
			log.Println("Attempt to register an existing user")
			return registerExist
		}
		// is the user allowed to log in?
		if u.state == 0 {
			log.Println("Attempt to log in unactivated account")
			return loginNotActive
		}
		// it's a valid login!

		// To keep the user logged in, we save two cookies.  The
		// first has the ID, the second has a random value.  We then
		// put ID->rand in our cookieStore map.  Subsequent requests
		// can grab the cookies and check the map
		//
		// NB: no Expires ==> it's a session cookie
		cookie := http.Cookie{Name: "id", Value: m["id"].(string), Path: "/"}
		http.SetCookie(w, &cookie)
		unique := sessionId()
		cookie = http.Cookie{Name: "key", Value: unique, Path: "/"}
		http.SetCookie(w, &cookie)
		cookieStore.Lock()
		cookieStore.m[m["id"].(string)] = unique
		cookieStore.Unlock()
		return loginOk
	}
}
