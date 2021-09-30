package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter().StrictSlash(true)

	// Add static folder
	staticDir := "/assets/"
	muxRouter.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	// Unauthorized
	muxRouter.HandleFunc("/unauthorized", unauthorized).Methods("GET")

	// Expired
	muxRouter.HandleFunc("/expired", expired).Methods("GET")

	// Logout
	muxRouter.HandleFunc("/logout/{login}", logout).Methods("GET")

	// Home private & public (public handle in func)
	muxRouter.HandleFunc("/", home).Methods("GET")
	muxRouter.HandleFunc("/", home).Methods("POST")

	// Encrypt messages
	muxRouter.HandleFunc("/encrypt/{verb}", encrypt).Methods("GET")
	muxRouter.HandleFunc("/encrypt", encryptPOST).Methods("POST")

	// Decrypt messages
	muxRouter.HandleFunc("/decrypt/{verb}", decrypt).Methods("GET")
	muxRouter.HandleFunc("/decrypt", decryptPOST).Methods("POST")

	// chat
	muxRouter.Handle("/chatSafe", authMiddleware(chatSafe)).Methods("GET")
	muxRouter.Handle("/chatSafe", authMiddleware(addMessage)).Methods("POST")

	// API
	muxRouter.Handle("/api/getMessages", authMiddleware(getMessages)).Methods("GET")

	return muxRouter
}
