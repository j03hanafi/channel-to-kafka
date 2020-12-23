package main

import "github.com/gorilla/mux"

func server() *mux.Router {
	router := mux.NewRouter()

	// endpoints
	router.HandleFunc("/payment/channel/iso", sendIso).Methods("POST")

	return router
}
