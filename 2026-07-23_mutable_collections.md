# Task: Mutable Collections
Date: 2026-07-23

## Goal
Implement mutable collections (`append`, `map_set`, `map_delete`) in the Zero transpiler. This is Improvement 31 in `improvements.md`.

## Steps
1. Examine `zero.go` to understand how lists and dicts are currently handled in AST generation.
2. Add support for `append`, `map_set`, and `map_delete` to the Go transpiler.
3. Update `orchestrator.py` to support the new grammar rules for these primitives.
4. Add tests in `test_mutable.zero` or similar.
5. Run `go build` and ensure everything compiles and tests pass.
6. Mark Improvement 31 as Done in `improvements.md`.
7. Groom backlogs as requested.
