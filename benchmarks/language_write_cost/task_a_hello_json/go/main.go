package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World! Zero language is alive!"))
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		msg := map[string]string{
			"status":  "success",
			"message": "Hello from Zero JSON endpoint!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	})

	http.ListenAndServe(":8080", nil)
}
