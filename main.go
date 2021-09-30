package main

import (
	"log"
)

func main() {
	// Customize log prefix
	setLog()

	// Start http server
	if web := runHTTP(); web != nil {
		log.Fatalf(web.Error())
	}
}
