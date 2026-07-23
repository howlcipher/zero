package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	content, err := os.ReadFile("names.txt")
	if err != nil {
		fmt.Println("Error: could not read names.txt")
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Println("Hello,", line)
		} else {
			fmt.Println("")
		}
	}
}
