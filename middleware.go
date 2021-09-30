package main

import (
	"log"
	"net/http"
	"time"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := getSession(r)
		if err != nil {
			log.Println("the session expired")
			log.Println(err)
			http.Redirect(w, r, "/expired", http.StatusTemporaryRedirect)
			return
		}
		now := time.Now().Add(time.Minute * 5).Unix() // Renew session 5 min before it ends

		if now < claims.ExpiresAt {
			if claims.Authorized && claims.Issuer == claimsIssuer {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/unauthorized", http.StatusForbidden)
			}
		} else {
			setSession(claims.Login, w)
			log.Printf("%s %s %s\n", "Session successfully renewed in", cookieName, "cookie")
			http.Redirect(w, r, "/chatSafe", http.StatusMovedPermanently)
		}

	})
}
