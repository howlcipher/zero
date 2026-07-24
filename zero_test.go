package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestOutputDirectoryFlag(t *testing.T) {
	// Build the zero binary
	cmd := exec.Command("go", "build", "-o", "zero", "zero.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build zero binary: %v", err)
	}
	defer os.Remove("zero")

	// Create a temporary directory for output
	outDir, err := os.MkdirTemp("", "zero-test-out-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outDir)

	// Create a dummy .zero file
	inputFile := filepath.Join(outDir, "dummy.zero")
	if err := os.WriteFile(inputFile, []byte("(cli_app (print \"Hello\"))"), 0644); err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}

	// Run the zero binary with -o flag
	cmd = exec.Command("./zero", "-o", outDir, inputFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run zero binary: %v", err)
	}

	// Check if server.go was created in the output directory
	serverFile := filepath.Join(outDir, "server.go")
	if _, err := os.Stat(serverFile); os.IsNotExist(err) {
		t.Errorf("Expected server.go to be created in %s, but it was not", outDir)
	}
}
