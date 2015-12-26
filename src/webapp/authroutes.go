// Routes related to OAuth, logging in, logging out, and registering
package main

import ("net/http"; "golang.org/x/oauth2"; "log")

// The route for '/login' starts the OAuth dance by redirecting to Google
//
// This is the first step of the OAuth dance.  Via oauthConf, we send
// ClientID, ClientSecret, and RedirectURL.  Step 2 is for google to check
// these fields and get the user to log in.
func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// NB: prefix the state string with 'l' so we can tell it's a login
	// later.  That's much easier than dealing with two different
	// redirect routes.
	url := oauthConf.AuthCodeURL("l"+oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// The route for '/register' is identical to '/login', except we change the
// state string to know it's a request to register.
func handleGoogleRegister(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL("r"+oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// The route for '/auth/google/callback' finishes the OAuth dance: It
// processes the Google response, and only sends us to the app page if
// the dance was successful.
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	res := processLoginReply(w, r)
	if res == loginOk {
		http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
		return
	}
	
	// we can't let the user in yet.  Set a flash cookie to explain, then send to '/'
	name := "eflash" // i or e for info or error
	val  := ""
	if res == registerOk {
		name = "iflash"
		val = "Registration succeeded.  Please wait for an administrator to confirm your account."
	} else if res == registerErr {
		val = "Registration error.  Please try again later."
	} else if res == registerExist {
		val = "The account you specified has already been registered."
	} else if res == loginNotReg {
		val = "The account you specified has not been registered."
	} else if res == loginNotActive {
		val = "The account you specified has not yet been activated by the administrator."
	} else if res == loginErr {
		val = "Login error.  Please try again later."
	} else {
		log.Fatal("processLoginReply() returned invalid status")
	}
	http.SetCookie(w, &http.Cookie{Name: name, Value: val, Path: "/"})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// The route for '/logout' logs the user out and redirects to home
func handleLogout(w http.ResponseWriter, r *http.Request) {
	processLogoutRequest(w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
