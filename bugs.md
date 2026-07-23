# 🐛 Bug Backlog

This document is the authoritative, ranked backlog for known flaws, bugs, and broken items specifically for the Zero transpiler project. It mirrors the structure of `improvements.md` and follows the same Working Protocol.

## Ranked Backlog (best ROI first)

Pending bugs carry the same diminishing-returns score defined in `improvements.md` (Score = Value × Decay ÷ Effort). Bugs rarely decay, so Decay is normally 1.0.

| # | Bug | Status | Score (V×D÷E) | Claude model | Gemini model | ROI rationale |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | [Lexer panics on EOF during unterminated string](#1-lexer-panics-on-eof-during-unterminated-string) | Done | 4.0 (4×1÷1) | Haiku 3 | Gemini 1.5 Flash | The lexer reports the error via `reportError`, but an explicit bounds check prevents potential runtime panics during deep parsing. |
| 2 | [Python Orchestrator timeout with heavy models](#2-python-orchestrator-timeout-with-heavy-models) | Done | 3.0 (3×1÷1) | Haiku 3 | Gemini 1.5 Flash | Added explicit UX warnings to console so users don't assume the script has frozen when loading heavy models in Outlines. |
| 3 | [Variable Shadowing in Go Generation](#3-variable-shadowing-in-go-generation) | Done | 4.0 (4×1÷1) | Sonnet 3.5 | Gemini 1.5 Pro | `let` expressions without inner `{}` brackets risk leaking scopes or redeclaring variables, causing Go compilation to fail on nested conditionals. |
| 4 | [Outlines EBNF Compilation Memory Limit](#4-outlines-ebnf-compilation-memory-limit) | Done | 2.5 (5×1÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Optimized the CFG to only enforce generic S-expression balanced parentheses, avoiding OOM on massive state machines while deferring semantic validation to the Go transpiler. |
| 5 | [AST Deep Nesting Stack Limits](#5-ast-deep-nesting-stack-limits) | Done | 2.0 (4×1÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Added an explicit depth parameter and limit of 1000 to `generateStatement` to fail gracefully with a JSON error instead of crashing the Go call stack. |
| 6 | [Include Paths Relative to File](#6-include-paths-relative-to-file) | Done (2026-07-22) | 6.0 (6×1÷1) | Sonnet 3.5 | Gemini 1.5 Pro | `(include "file.zero")` uses current working directory rather than resolving paths relative to the file doing the inclusion. |
| 7 | [defun Typing Rigidness](#7-defun-typing-rigidness) | Done (2026-07-22) | 1.75 (7×1÷4) | Sonnet 3.5 | Gemini 1.5 Pro | All `defun` arguments strictly compile to Go `string` types, breaking if we try to pass an `*http.Request` or `int` to a function. |
| 8 | [try_let Error Interception Rigidness](#8-try_let-error-interception-rigidness) | Done (2026-07-23) | 1.66 (5×1÷3) | Sonnet 3.5 | Gemini 1.5 Pro | `try_let` is currently hardcoded to only support `parse_json` as the error-returning function. Needs generalization. |
| 11 | [No Runtime Source Mapping](#11-no-runtime-source-mapping) | Done (2026-07-23) | 6.0 (6×1÷1) | Haiku 3 | Gemini 1.5 Flash | Go panics at runtime do not map back to `.zero` files. Need Go `//line` directives. |
| 12 | [Lexer cannot tokenize `!=`](#12-lexer-cannot-tokenize-) | Pending | 8.0 (8×1÷1) | Haiku 3 | Gemini 1.5 Flash | The lexer's symbol character class omits `!`, so any script using `!=` fails to lex at all — including the README's own `test` block example. One-line fix. |
| 13 | [`return` in `defun` drops compound expressions](#13-return-in-defun-drops-compound-expressions) | Pending | 3.5 (7×1÷2) | Sonnet 3.5 | Gemini 1.5 Pro | `(return (+ a b))` or `(return (call f x))` emits a bare `return` with no value, which fails Go compilation for any non-void `defun`. Only bare symbols/literals currently work. |
| 14 | [`(import)` duplicated/unused in `server_test.go`](#14-import-pkg-duplicatedunused-in-generated-server_testgo) | Pending | 3.0 (6×1÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Custom imports get blindly copied into both `server.go` and `server_test.go`, causing duplicate-import or unused-import compile failures once a script combines `import`, `defun`, and `test` blocks. |
| 10 | [String Escaping Limitations](#10-string-escaping-limitations) | Done | 5.0 (5×1÷1) | Haiku 3 | Gemini 1.5 Flash | Lexer breaks on unicode escapes and escaped single quotes, common in LLM outputs. |
| 9 | [Depth Limit Crash via `let` Chaining](#9-depth-limit-crash-via-let-chaining) | Done (2026-07-23) | 4.0 (4×1÷1) | Haiku 3 | Gemini 1.5 Flash | Long sequential scripts with variable assignments crash transpiler with `AST too deep`. Scope needs flattening. |
## Details

### 1. Lexer panics on EOF during unterminated string
In `zero.go`, the string lexer currently scans for the next quote. If EOF is hit, it correctly calls `reportError`, but if the file is truly truncated, we should ensure no other routines attempt to consume past the array bounds. The code is mostly safe but needs explicit unit tests to guarantee it never panics outside of the controlled JSON error output.

### 2. Python Orchestrator timeout with heavy models
When running `orchestrator.py` against a local Ollama instance with a large model (e.g. 70B parameters), the structured generation engine in `outlines` might hang for a long time compiling the regex/grammar. We need to add logging or a loading indicator so the user knows the compilation/generation hasn't silently frozen.

### 6. Include Paths Relative to File
* **Description:** `(include "file.zero")` uses `os.ReadFile(filename)`.
* **Why:** If the zero compiler is run from a different directory, it will fail to find relative included files.
* **Impact:** 6/10 (High).

### 7. defun Typing Rigidness
* **Description:** `defun` maps all arguments to strings in Go.
* **Why:** We cannot cleanly pass `req` to helper functions without using strings.
* **Impact:** 7/10 (High).

### 8. try_let Error Interception Rigidness
* **Description:** `try_let` only intercepts `parse_json` right now.
* **Why:** Needs to generalize to any function returning `(value, error)`.
* **Impact:** 5/10 (Medium).

### 9. Depth Limit Crash via `let` Chaining
* **Description:** Every `let` block increments the AST depth (`depth+1` at line 644). If an AI generates a long sequential script with 1,000 variable assignments, the transpiler crashes with `AST too deep`.
* **Why:** Variable scope needs to be flattened instead of infinitely nesting in AST.
* **Impact:** 8/10 (High).

### 10. String Escaping Limitations
* **Description:** The lexer only supports `\n` and `\t`. It breaks on unicode escapes (`\uXXXX`) and escaped single quotes (`\'`).
* **Why:** LLMs frequently output these escapes, especially in JSON strings.
* **Impact:** 7/10 (High).

### 11. No Runtime Source Mapping
* **Description:** Runtime panics point to `main.go`, not the `.zero` file.
* **Why:** The transpiler should inject Go `//line filename.zero:line` directives so runtime panics map back to the AI's code.
* **Impact:** 6/10 (Medium).

### 12. Lexer cannot tokenize `!=`
* **Description:** `lexer()` in `zero.go` accepts `=`, `<`, `>`, `+`, `-`, `*`, `.`, `/`, `_` as symbol start/continuation characters (~line 140), but not `!`. Any `.zero` source containing `!=` fails at the lex stage with `Unexpected character: !`, before the parser is ever reached — even though `generateStatement`'s `if`/`match` condition handling (lines ~830, ~930) explicitly checks for and supports an `"!="` operator string. The operator support is dead code today.
* **Why:** Not-equal comparisons are basic control flow. Discovered while verifying improvement #16 (Native Unit Test Blocks): the exact `(if (!= result 5) ...)` example added to the README's "Native Unit Test Blocks" section (and `docs/index.html`) does not lex, so the flagship new-feature example is broken as documented.
* **Impact:** 8/10 (High — silently breaks a documented, common comparison operator).
* **Repro:** `echo '(cli_app (if (!= 1 2) (print "ne")))' > /tmp/t.zero && ./zero /tmp/t.zero` → `{"reason":"Unexpected character: !","line":1,...}`.
* **Fix sketch:** add `ch == '!'` to both character-class checks at zero.go:140 and zero.go:142.

### 13. `return` in `defun` drops compound expressions
* **Description:** In `generateStatement`, the `"return"` case (zero.go:597-606) does `fmt.Sprintf("return %s", valNode.Value)`. `Node.Value` is only populated for leaf tokens (SYMBOL/STRING/NUMBER) — for a compound expression node like `(+ a b)` or `(call f x)`, `Value` is empty, so the case silently emits a bare `return` with nothing after it. Every other statement/expression case recursively calls `generateStatement(valNode, reqVar, depth+1)`; `return` is the outlier.
* **Why:** Discovered via the same repro as bug #12 — even after fixing `!=` lexing, the README's `(defun add (a b) (return (+ a b)))` example fails to compile in the generated Go (`not enough return values`), because the return value is silently dropped. This affects any `defun` whose return value is not a bare variable or literal (arithmetic, `call`, string concatenation, etc.), a large fraction of realistic function bodies.
* **Why it went unnoticed:** the transpiler exits 0 and produces syntactically-valid-looking Go for simple `(return x)`/`(return "str")` cases; the failure only shows up as a downstream `go build` error on the generated code, which existing example `.zero` files apparently don't exercise.
* **Impact:** 7/10 (High — undermines improvement #6 Function Definitions for any non-trivial return value).
* **Repro:** `(cli_app (defun add (a b) (return (+ a b))) (print (call add 2 3)))` → generated `func add(a string, b string) string { return \n}`.
* **Fix sketch:** replace `valNode.Value` with `generateStatement(valNode, reqVar, depth+1)` in the non-STRING branch, mirroring how other cases (e.g. `call`, binary ops) already recurse.

### 14. `(import "pkg")` duplicated/unused in generated `server_test.go`
* **Description:** `generateCode` appends every `extraImports` entry into both `server.go`'s and `server_test.go`'s import blocks unconditionally (zero.go:488 and zero.go:540). Two problems: (a) if a custom import happens to name a package already in the hardcoded default import list (e.g. `(import "strings")`), the generated file declares it twice → Go "imported and not used" / redeclare compile error; (b) even for non-colliding packages, `server_test.go` has no `var _ = pkg.Symbol` suppression the way the default imports do, so a package used only by non-test code (routes/defuns) but not referenced inside any `(test ...)` body makes `server_test.go` fail with "imported and not used" as soon as any `test` block exists alongside any `import`.
* **Why:** Found while verifying improvement #16 (Native Unit Test Blocks) end-to-end with a `.zero` file that combines `(import ...)` with a `(defun ...)` and a `(test ...)` block — a realistic combination now that both features exist.
* **Impact:** 6/10 (Medium-High — breaks the common case of testing a `defun` that relies on an external package).
* **Repro:** `(cli_app (import "strings") (defun shout (s) (return (call strings.ToUpper s))) (test "d" (print (call shout "hi"))))` → `server.go` and `server_test.go` both declare `"strings"` twice in their import blocks.
* **Fix sketch:** dedupe `extraImports` against the hardcoded default list before appending; for `server_test.go`, either only include imports actually referenced within test bodies, or emit `var _ = ` blanks for extraImports the same way defaults are handled.
