package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// Set CORS
func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Run server
func runHTTP() error {

	httpStart := "http server listening on port:"
	httpFailed := "http server failed to listening on port:"

	mux := makeMuxRouter()

	httpPort := getEnv("HTTPPORT", "8888")

	log.Printf("%s %s \n", httpStart, httpPort)

	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        cors(mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return errors.New(httpFailed)

	}
	return nil
}
