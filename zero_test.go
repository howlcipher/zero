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

func TestCrashStateSerialization(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "zero", "zero.go")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build zero binary: %v", err)
	}
	defer os.Remove("zero")

	outDir, err := os.MkdirTemp("", "zero-crash-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outDir)

	inputFile := filepath.Join(outDir, "panic.zero")
	if err := os.WriteFile(inputFile, []byte(`(cli_app (call panic "test crash dump"))`), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	cmd = exec.Command("./zero", "-o", outDir, inputFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to transpilation: %v", err)
	}

	serverFile := filepath.Join(outDir, "server.go")
	appBin := filepath.Join(outDir, "app")
	cmd = exec.Command("go", "build", "-o", appBin, serverFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build generated server.go: %v", err)
	}

	cmd = exec.Command(appBin)
	cmd.Dir = outDir
	_ = cmd.Run() // expected to exit with non-zero code

	crashFile := filepath.Join(outDir, "crash.json")
	data, err := os.ReadFile(crashFile)
	if err != nil {
		t.Fatalf("Failed to read crash.json: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("crash.json is empty")
	}
}

