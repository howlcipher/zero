# Improvement 48: Add CLI flag for output directory

## Overview
Implemented the `-o` CLI flag for the `zero` transpiler to specify the output directory for `server.go` and `server_test.go`.

## Changes Made
- Modified `zero.go` to import the `flag` package.
- Added `outDir := flag.String("o", "", "output directory")` to parse the CLI flag.
- Updated the file writing logic to use `filepath.Join(*outDir, ...)` for `server.go` and `server_test.go`.
- Added a Go test `zero_test.go` to verify the functionality of the `-o` flag.
- Updated `improvements.md` to mark item 48 as Done.
