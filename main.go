package main

import (
	"log"
	"net/http"
	"sona/service"

	"sona/gen/sonav1connect"
)

func main() {
	// Create the server instances
	helloServer := service.NewHelloServer()
	mux := http.NewServeMux()

	// Mount the ConnectRPC handlers
	helloPath, helloHandler := sonav1connect.NewHelloServiceHandler(helloServer)
	mux.Handle(helloPath, helloHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
