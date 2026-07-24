# Task Journal: Improvement #20 Auto-Tracing

## Overview
Implemented the `(trace var)` macro to natively support debugging without manually formatting `fmt.Println` statements.

## Changes Made
- Modified `zero.go` to add `"trace"` keyword in `generateStatementRaw`.
- `trace` compiles to `fmt.Println("[filename.zero:line] var =", var)` matching the requested format.
- `AI_PROMPT.md` was updated to document the `trace` macro.
- Verified that `orchestrator.py` does not need changes since its context-free grammar already supports generic symbols, and `trace` fits exactly into the existing generic S-expression matcher.
- Added `tests/test_trace.zero` and verified it compiles and runs correctly.
- Groomed `improvements.md` to mark Improvement 20 as Done.

## Testing
- Used `go build -o zero zero.go` to compile the transpiler.
- Executed `./zero tests/test_trace.zero && go run server.go` and verified the output format matches `[tests/test_trace.zero:3] a = 42`.
