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
		fmt.Println("Hello, World!")
		{
			name := "Zero"
			_ = name
		fmt.Println("Welcome to", name)
		}
}
