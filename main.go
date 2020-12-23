package main

import (
	"log"
	"net/http"
)

func main() {

	router := server()
	log.Fatal(http.ListenAndServe(":6010", router))

}
