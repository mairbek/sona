package main

import (
	"log"
	"net/http"
	"sona/service"

	"sona/gen/sonav1connect"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Create the server instances
	helloServer := service.NewHelloServer()
	mux := http.NewServeMux()

	// Mount the ConnectRPC handlers
	helloPath, helloHandler := sonav1connect.NewHelloServiceHandler(helloServer)
	mux.Handle(helloPath, helloHandler)

	// Legacy router
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/api/v1/pepe", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pepuuu"))
	})

	mux.Handle("/api/", router)

	// Create and mount the path rewrite interceptor
	// Define path rewrite rules
	pathRewriteRules := []service.PathRewriteRule{
		{
			From: "/api/v1/hello",
			To:   "/sona.v1.HelloService/Hello",
		},
		// Add more rules here as needed
	}
	interceptor := service.NewPathRewriteInterceptor(pathRewriteRules, mux)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", interceptor); err != nil {
		log.Fatal(err)
	}
}
