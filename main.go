package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	// logging
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal("Found error in log ", err)
	}

	log.SetOutput(file)

	// setup server
	router := server()
	log.Fatal(http.ListenAndServe(":6010", router))

}
