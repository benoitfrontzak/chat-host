package main

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Generate JSON Web Token
// A token is made of three parts, separated by .'s.
// The first part is called the header.
// It contains the necessary information for verifying the last part, the signature.
// The first two parts are JSON objects, that have been base64url encoded.
// The part in the middle is the interesting bit. It's called the Claims and contains the actual stuff you care about.
// Refer to RFC 7519 for information about reserved keys and the proper way to add your own.
// The last part is the signature, encoded the same way.

func generateJWT(login string, exp time.Time) (string, error) {
	var ck = []byte(cipherKey)
	tk := jwt.New(jwt.SigningMethodHS256)
	c := tk.Claims.(jwt.MapClaims)
	e := exp.Unix()

	c["authorized"] = true
	c["login"] = login
	c["iss"] = claimsIssuer // issuer
	c["exp"] = e

	tkStr, err := tk.SignedString(ck)
	if err != nil {
		return "", err
	}

	return tkStr, nil
}

// Custom claims for JSON Web Token
// type StandardClaims struct {
// 	Audience  string `json:"aud,omitempty"`
// 	ExpiresAt int64  `json:"exp,omitempty"`
// 	Id        string `json:"jti,omitempty"`
// 	IssuedAt  int64  `json:"iat,omitempty"`
// 	Issuer    string `json:"iss,omitempty"`
// 	NotBefore int64  `json:"nbf,omitempty"`
// 	Subject   string `json:"sub,omitempty"`
// }
type customClaims struct {
	Login      string `json:"login"`
	Authorized bool   `json:"authorized"`
	jwt.StandardClaims
}

func returnClaims(tk string) *customClaims {
	claims := &customClaims{}
	tkn, err := jwt.ParseWithClaims(tk, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cipherKey), nil
	})
	if err != nil {
		log.Println(err)
	}
	if !tkn.Valid {
		return nil
	}
	return claims
}
