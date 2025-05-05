package main

import (
	"log"
	"net/http"
	"sona/service"

	"sona/gen/sonav1connect"

	"github.com/rs/cors"
)

func main() {
	// Create the server instances
	helloServer := service.NewHelloServer()
	mux := http.NewServeMux()

	// Define path rewrite rules
	pathRewriteRules := []service.PathRewriteRule{
		{
			From: "/v1/hello",
			To:   "/sona.v1.HelloService/Hello",
		},
		// Add more rules here as needed
	}

	// Mount the ConnectRPC handlers
	helloPath, helloHandler := sonav1connect.NewHelloServiceHandler(helloServer)
	mux.Handle(helloPath, helloHandler)

	// Create and mount the path rewrite interceptor
	interceptor := service.NewPathRewriteInterceptor(pathRewriteRules, mux)

	// Add CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}).Handler(interceptor)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatal(err)
	}
}
