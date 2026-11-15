package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":9060"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request for path: %s\n", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello from the backend server!"}`))
	})

	fmt.Printf("Mock server is running on http://localhost%s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
