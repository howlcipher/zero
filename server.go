package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)
func main() {
	var _ = sql.Open
	var _ = os.Getenv
	var _ = json.Marshal
	var _ = io.ReadAll
	var _ = http.DefaultClient
	http.HandleFunc("/err", func(w http.ResponseWriter, req *http.Request) {
//line err.zero:3
		{
			x := (1 + "a")
			_ = x
//line err.zero:4
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		fmt.Fprint(w, "OK")
		}
	})
	
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
