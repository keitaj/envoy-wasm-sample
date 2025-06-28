package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "Hello from backend!",
			"headers": map[string]string{
				"x-auth-user":  r.Header.Get("x-auth-user"),
				"x-request-id": r.Header.Get("x-request-id"),
			},
			"path":   r.URL.Path,
			"method": r.Method,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Backend service starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
