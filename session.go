package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

// securecookie encodes and decodes authenticated and optionally encrypted cookie values.
// Secure cookies can't be forged, because their values are validated using HMAC.
// When encrypted, the content is also inaccessible to malicious eyes.
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func setSession(login string, w http.ResponseWriter) {
	// Define 'expires' and 'maxAge' for cookie
	// and jwt claims 'expiredAt' (unix time of expiresDate)
	expiresDate := time.Now().Add(time.Minute * 60) // 60 min
	maxAge := 60 * 60                               // 60 min
	// Generate JSON Web Token add login/role to claims
	jwt, err := generateJWT(login, expiresDate)
	if err != nil {
		log.Println(err)
	}
	// Secure cookie
	// Encode token with cookie name as value of the cookie
	if encoded, err := cookieHandler.Encode(cookieName, jwt); err == nil {
		cookie := &http.Cookie{
			Name:       cookieName,
			Value:      encoded,
			Path:       "/",
			Domain:     "",
			Expires:    expiresDate,
			RawExpires: "",
			MaxAge:     maxAge,
			Secure:     false,
			HttpOnly:   true,
			SameSite:   http.SameSiteStrictMode,
			Raw:        "",
			Unparsed:   []string{},
		}
		http.SetCookie(w, cookie)
	}
}

// Returns the claims of the http request
func getSession(r *http.Request) (*customClaims, error) {
	var claims *customClaims

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	var cookieValue string
	if err = cookieHandler.Decode(cookieName, cookie.Value, &cookieValue); err == nil {
		claims = returnClaims(cookieValue)
	}

	return claims, nil
}

// Clear the cookie
func clearSession(w http.ResponseWriter) {
	log.Printf("%s %s %s \n", "the session", cookieName, "have been deleted successfully")
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
