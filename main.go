package main

import (
	"log"
	"net/http"
	"sona/service"

	"sona/gen/sonav1connect"

	"github.com/rs/cors"
)

func main() {
	helloServer := service.NewHelloServer()
	mux := http.NewServeMux()

	// Mount the ConnectRPC handler
	path, handler := sonav1connect.NewHelloServiceHandler(helloServer)
	mux.Handle(path, handler)

	// Add CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}).Handler(mux)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatal(err)
	}
}
