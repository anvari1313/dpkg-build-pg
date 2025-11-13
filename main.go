package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Register the handler for the root path
	http.HandleFunc("/", handleRoot)

	// Start the server on port 8080
	port := ":8080"
	log.Printf("Starting HTTP server on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleRoot is the handler function that responds with a string
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from dpkg-build-pg server!")
}
