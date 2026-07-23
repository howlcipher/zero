# 🚀 Improvement Backlog

This document is the authoritative, ranked backlog for the Zero transpiler project. It mirrors the format used in the main AI Knowledge Library.

## Working Protocol

This protocol applies to every worked task in the Zero project:

1. **Open a task journal.** Record your steps in a `YYYY-MM-DD_task_name.md` file if the task is complex.
2. **Re-evaluate the model.** Pick the least expensive available model (e.g., local Ollama, Claude, or Gemini) that can do the job well for the Zero transpiler.
3. **Route the crafted skills.** Check `.agents/skills/zero_transpiler/SKILL.md` before planning.
4. **Scan for helpful free tools.** Ensure you aren't rebuilding something already available.
5. **Finish the loop.** Every code change ships with relevant tests. Run Go builds (`go build`) and Python script validations before committing.

## Ranked Backlog (best ROI first)

Pending rows are ranked by a diminishing-returns score:

**Score = (Value × Decay) ÷ Effort**
- **Value (1–8):** pain or risk removed if the item ships.
- **Decay:** geometric halving per already-shipped item in the same theme (1.0 → 0.5 → 0.25 …).
- **Effort (1–8):** roughly log-scale; 1 = minutes, 8 = weeks.

| # | Improvement | Status | Score (V×D÷E) | Claude model | Gemini model | ROI rationale |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | [Add Routing Support](#1-add-routing-support) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Highest value to allow building web apps with multiple endpoints instead of just the root path. |
| 3 | [Extend Python Orchestrator Grammar](#3-extend-python-orchestrator-grammar) | Done | — | Haiku 3 | Gemini 1.5 Flash | Must update the grammar in `orchestrator.py` immediately after adding new Go AST features so the LLM can use them. |
| 2 | [Add Conditionals and Variables](#2-add-conditionals-and-variables) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Necessary for basic logic flow in handlers (checking methods, parsing headers). |
| 4 | [Add Database Connections (SQL)](#4-add-database-connections-sql) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Crucial for dynamic data and actual web service capabilities. |
| 5 | [Add JSON Request/Response Handling](#5-add-json-requestresponse-handling) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Needed to build standard REST APIs. Decay 0.125 because three Go AST features shipped. |
| 6 | [Add Function Definitions (defun)](#6-add-function-definitions-defun) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Critical for code modularity (DRY principle). |
| 7 | [Add Structs and Type Definitions](#7-add-structs-and-type-definitions) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Necessary for strict Input Validation schemas, adhering to software_development skill guidelines, and mapping SQL/JSON to Go. |
| 8 | [Add Iteration and Data Structures](#8-add-iteration-and-data-structures) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Essential for handling arrays of SQL results (list, map, for). |
| 9 | [Add Environment Variables Access](#9-add-environment-variables-access) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Follows 'Secure by Default' guidelines to prevent hardcoding database credentials or secrets in S-expressions. Decay 0.125. |
| 10 | [Add External Module Imports](#10-add-external-module-imports) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Allows importing third-party Go packages, unlocking the entire Go ecosystem. Decay 0.125. |
| 11 | [Add Concurrency (spawn)](#11-add-concurrency-spawn) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Allows AI to effortlessly run background jobs without blocking HTTP responses. |
| 12 | [Add Error Handling (try/catch)](#12-add-error-handling-trycatch) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Crucial for safe execution. Maps to Go's `if err != nil` idiom. |
| 13 | [Add File Inclusions (include)](#13-add-file-inclusions-include) | Done | 2.33 (7×1.0÷3) | Sonnet 3.5 | Gemini 1.5 Pro | Prevents massive monolithic `.zero` files by allowing modular codebases. |
| 14 | [Add Basic Math and Logic Operators](#14-add-basic-math-and-logic-operators) | Done | — | Sonnet 3.5 | Gemini 1.5 Pro | Necessary for computing values natively in Zero instead of relying entirely on DB logic. |
| 15 | [Add Middleware Support](#15-add-middleware-support) | Done | 0.41 (5×0.25÷3) | Sonnet 3.5 | Gemini 1.5 Pro | Required for adding authentication and request logging across routes. |
| 35 | [Add LLM-powered Type Coercion (fuzzy_cast)](#35-add-llm-powered-type-coercion-fuzzy_cast) | Pending | 2.00 (8×1.0÷4) | Sonnet 3.5 | Gemini 1.5 Pro | Universal parser using LLM structured outputs to map messy, unstructured text to strict structs. |
| 36 | [Add Intent-based Validation (assert_semantic)](#36-add-intent-based-validation-assert_semantic) | Pending | 2.00 (6×1.0÷3) | Sonnet 3.5 | Gemini 1.5 Pro | Enforces complex, qualitative natural language boundaries effortlessly using zero-shot prompts. |
| 34 | [Add Semantic Routing (semantic_match)](#34-add-semantic-routing-semantic_match) | Pending | 1.40 (7×1.0÷5) | Sonnet 3.5 | Gemini 1.5 Pro | Natively understands intent, replacing brittle traditional conditional routing and regexes. |
| 37 | [Add Just-In-Time Function Generation (lazy_synthesize)](#37-add-just-in-time-function-generation-lazy_synthesize) | Pending | 0.71 (5×1.0÷7) | Sonnet 3.5 | Gemini 1.5 Pro | Defers boilerplate generation to runtime, allowing AI to focus only on high-level logic. |
| 41 | [Add Stochastic Control Flow](#41-add-stochastic-control-flow) | Pending | 0.29 (2×1.0÷7) | Sonnet 3.5 | Gemini 1.5 Pro | Introduces fuzzy logic natively; deferred as non-essential for initial MVP. |
| 38 | [Add Swarm Primitives](#38-add-swarm-primitives) | Pending | 0.25 (2×1.0÷8) | Sonnet 3.5 | Gemini 1.5 Pro | Extremely advanced futurist concept; deferred to maintain MVP scope. |
| 39 | [Add Teleological Execution](#39-add-teleological-execution) | Pending | 0.25 (2×1.0÷8) | Sonnet 3.5 | Gemini 1.5 Pro | Radical paradigm shift, non-critical enhancement deferred from MVP. |
| 40 | [Add Auto-Mutating Runtime](#40-add-auto-mutating-runtime) | Pending | 0.12 (1×1.0÷8) | Sonnet 3.5 | Gemini 1.5 Pro | Highly experimental runtime evolution; deferred per strict MVP boundaries. |

## Details

### 1. Add Routing Support
* **Description:** Update the compiler to accept multiple `(route path handler)` definitions inside a web server block.
* **Why:** The prototype only builds a single server with a hardcoded route. Real applications need routers.
* **Impact:** 2/10 (Minor - helpful but not strictly blocking).

### 2. Add Conditionals and Variables
* **Description:** Introduce `let` and `if` blocks to handle internal request logic. For example: `(if (= req.method "POST") ...)`. This will require updating the Lexer to handle operators like `=` and the Code Generator to output Go `if` statements.
* **Why:** Web handlers need to implement dynamic logic based on request types and data.
* **Impact:** 8/10 (High).

### 3. Extend Python Orchestrator Grammar
* **Description:** Currently, `orchestrator.py` uses a strict regex for the proof-of-concept single endpoint. As we implement improvements 1 and 2, this regex needs to be translated into a full Context Free Grammar (CFG) using Outlines to support nested expressions and arbitrary routes.
* **Why:** The LLM agent loop breaks if it cannot generate valid syntax for new AST nodes.
* **Impact:** 4/10 (Medium - blocks orchestrator but not manual transpiler usage).

### 4. Add Database Connections (SQL)
* **Description:** Implement SQL database connections via Go's `database/sql` mapping to an S-expression like `(sql_query db "SELECT * FROM users")`.
* **Why:** Real-world applications require state and persistence.
* **Impact:** 6/10 (Medium).

### 5. Add JSON Request/Response Handling
* **Description:** Implement a way to parse JSON bodies into variables and output JSON responses cleanly via `encoding/json`. E.g., `(parse_json req.body)` and `(res_json 200 data)`.
* **Why:** The modern web runs on JSON; text/plain is insufficient.
* **Impact:** 5/10 (Medium).

### 6. Add Function Definitions (defun)
* **Description:** Allow defining standard functions `(defun name (args) body)` outside of routes that can be called anywhere.
* **Why:** Needed to adhere to modularity and DRY principles.
* **Impact:** 8/10 (High).

### 7. Add Structs and Type Definitions
* **Description:** Implement `(struct Name (field type) ...)` to enforce Go's strict typing system for parsing JSON and scanning SQL rows.
* **Why:** Strictly typed inputs are a core requirement of defensive programming and input validation.
* **Impact:** 7/10 (High).

### 8. Add Iteration and Data Structures
* **Description:** Support loops `(for ...)` and basic collections `(list ...)` and `(dict ...)`.
* **Why:** Essential for mapping over database query results or iterating through JSON arrays.
* **Impact:** 6/10 (Medium).

### 9. Add Environment Variables Access
* **Description:** Introduce a `(env "KEY")` node to retrieve environment variables.
* **Why:** Vital for securely injecting database credentials and API keys without hardcoding them in the S-expressions.
* **Impact:** 3/10 (Low/Medium - security critical).

### 10. Add External Module Imports
* **Description:** Allow defining `(import "github.com/pkg")` at the root level to pull in external Go code.
* **Why:** Makes Zero extensible and leverages the massive open-source Go ecosystem.
* **Impact:** 3/10 (Low - advanced feature).

### 11. Add Concurrency (spawn)
* **Description:** Add a `(spawn (lambda () ...))` node that maps to Go's `go func() {}` to execute non-blocking routines.
* **Why:** AI agents building web applications often need to trigger background processes (like sending emails or metrics) without delaying the HTTP response.
* **Impact:** 7/10 (High).

### 12. Add Error Handling (try/catch)
* **Description:** Implement `(try (expression) (catch err ...))` to wrap Go expressions that return `(value, error)`. 
* **Why:** Go relies heavily on `if err != nil`. We need a clean, Lisp-like way to handle these errors safely in Zero without panicking.
* **Impact:** 8/10 (High - critical for production safety).

### 13. Add File Inclusions (include)
* **Description:** Implement `(include "routes.zero")` to dynamically merge multiple Zero files during the transpilation step.
* **Why:** A full-fledged language needs modularity. Right now, everything must live in one massive S-expression.
* **Impact:** 7/10 (High).

### 14. Add Basic Math and Logic Operators
* **Description:** Support native mathematical and logical operators like `(+ 1 2)`, `(- a b)`, `(and x y)`.
* **Why:** Computing logic natively (like paginating data or computing totals) is currently impossible without external SQL/Go functions.
* **Impact:** 8/10 (High).

### 15. Add Middleware Support
* **Description:** Introduce a `(middleware auth_func)` block that can wrap a set of `(route ...)` blocks.
* **Why:** Modern APIs require authentication headers, logging, and CORS handling. Middleware is the standard pattern for this.
* **Impact:** 5/10 (Medium).

### 34. Add Semantic Routing (semantic_match)
* **Description:** A control flow structure that routes execution based on the semantic proximity (intent and meaning) of an input string compared to a set of natural language descriptions.
* **Why:** Natively understands intent. Acknowledges that human language is fuzzy and allows the code to handle it gracefully without exhaustive mapping or complex regexes.
* **Impact:** 7/10 (High - unlocks intent-based routing).

### 35. Add LLM-powered Type Coercion (fuzzy_cast)
* **Description:** A casting function `fuzzy_cast[T]` that uses structured-output LLM APIs to automatically coerce messy, unstructured text into a strictly typed struct `T`.
* **Why:** Traditional serialization requires perfect 1:1 schema matches. This acts as a universal, intelligent parser that infers required mapping.
* **Impact:** 8/10 (High - eliminates brittle parsing code).

### 36. Add Intent-based Validation (assert_semantic)
* **Description:** An assertion primitive that evaluates qualitative, subjective natural language conditions against a variable. E.g. `assert_semantic(user_bio, "is professional")`.
* **Why:** Allows the code to enforce complex, qualitative boundaries effortlessly, removing the need for massive heuristic functions.
* **Impact:** 6/10 (Medium - powerful for data safety).

### 37. Add Just-In-Time Function Generation (lazy_synthesize)
* **Description:** A declarative primitive for defining a function using only its signature and a natural language docstring. The implementation is dynamically generated the first time it is invoked.
* **Why:** AI writing the language doesn't have to waste tokens generating mundane utility functions, delegating implementation to the runtime.
* **Impact:** 5/10 (Medium - innovative but complex to execute).

### 38. Add Swarm Primitives
* **Description:** Introduces autonomous subagents as first-class concurrency objects. Developers orchestrate a swarm of agents using primitives like `(spawn_agent "Researcher" (task "find sources"))` that communicate via typed message-passing channels and autonomously negotiate tasks.
* **Why:** Concurrency shifts from deterministic CPU scheduling to non-deterministic, autonomous orchestration, breaking conventional rules and allowing agents to independently verify upstream outputs.
* **Impact:** 2/10 (Low - extremely advanced, deferred for strict MVP scoping).

### 39. Add Teleological Execution
* **Description:** A goal-driven syntax where developers define a target state (e.g., `(achieve (is_sorted list) (using "quick sort algorithm"))`) rather than imperative steps. The runtime acts as a solver to dynamically search for the execution path and execute necessary steps.
* **Why:** Abandons imperative control flow entirely. Code becomes a set of constraints and objectives, making execution a continuous planning and state-space search process.
* **Impact:** 2/10 (Low - radical shift, deferred for MVP).

### 40. Add Auto-Mutating Runtime
* **Description:** A self-rewriting primitive `(optimize_block ...)` that monitors execution metrics and automatically employs an LLM to rewrite and hot-swap its underlying Go implementation at runtime if bottlenecks are detected.
* **Why:** Code becomes active and evolutionary in production rather than immutable, natively incorporating model evaluation and code generation into the execution cycle.
* **Impact:** 1/10 (Low - highly experimental).

### 41. Add Stochastic Control Flow
* **Description:** Natively handles uncertainty in the AST. Conditions evaluate to probability distributions, allowing control flow primitives like `(if (> (confidence (is_fraud tx)) 0.95) ...)` to branch based on statistical certainty.
* **Why:** Eliminates hardcoded heuristics by bringing fuzzy logic directly into the core execution loop, perfectly matching the probabilistic nature of AI models.
* **Impact:** 3/10 (Low/Medium - complex but powerful for AI).

---

## V2: AI-First Language Optimizations

Now that Zero V1 is complete (a full Turing-complete web server and CLI language), the next phase is optimizing it specifically for **Autonomous AI Development**. Since Zero does not need to be human-readable, we can bend the language features to perfectly suit AI agents.

### Proposed Improvements

| # | Improvement | Status | Score | AI Rationale |
| --- | --- | --- | --- | --- |
| 17 | **Type Hinting for `defun`** | Done (2026-07-22) | 3.5 (7×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Currently, all `defun` arguments compile to `string`. Adding `(type_hint var "int")` ensures the AI gets immediate compile-time errors from Go. |
| 19 | **Context/Intent Nodes (`intent`)** | Done (2026-07-22) | 2.0 (4×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | `(intent "I am building a login flow")`. The transpiler strips these out, but agents can parse them to instantly understand context. |
| 21 | **Native HTTP Client (`fetch`)** | Done (2026-07-23) | 4.0 (8×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Essential for an AI language to interact with external APIs (like LLM providers or GitHub) without writing raw Go `net/http` code. |
| 22 | **Subprocess Execution (`exec`)** | Pending | 3.5 (7×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Crucial for automation tasks (e.g. `(exec "git status")`). Follows automation skills for script consolidation. |
| 23 | **File I/O Operations (`read_file`)** | Pending | 3.0 (6×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Needed to replace Bash/Python for file manipulation. `(write_file "log.txt" data)` and `(read_file "config.json")`. |
| 24 | **CLI Argument Parsing (`cli_args`)** | Pending | 2.5 (5×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Required for workflow consolidation (per `automation` skill). Allows Zero scripts to take parameters effortlessly. |
| 25 | **Timers and Backoff (`sleep`)** | Pending | 2.0 (4×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | Fault tolerance (per `automation` skill) requires exponential backoff and deliberate delays `(sleep 1000)` during API rate limits. |
| 16 | **Native Unit Test Blocks (`test`)** | Pending | 1.5 (6×1.0÷4) | Sonnet 3.5 | Gemini 1.5 Pro | AI iterates faster with TDD. A native `(test "desc" ...)` block at the root that compiles directly to `_test.go` allows seamless testing. |
| 20 | **Auto-Tracing (`trace`)** | Pending | 1.5 (3×1.0÷2) | Sonnet 3.5 | Gemini 1.5 Pro | AI debugs by spamming `print`. A `(trace var)` macro auto-injects line numbers and variable names into `fmt.Println`. |
| 18 | **Declarative Schema Migrations** | Pending | 1.0 (5×1.0÷5) | Sonnet 3.5 | Gemini 1.5 Pro | If `(schema "users" (column "id" "int"))` is in `.zero`, the transpiler can auto-generate `CREATE TABLE IF NOT EXISTS`. |
| 26 | **LLM-Native Primitives (`llm_generate`)** | Pending | 6.0 (6×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | Built-in nodes like `(llm_generate "prompt")` to make it trivial for an AI to utilize other AIs. |
| 27 | **AST-Level Semantic Patching** | Pending | 5.0 (5×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | `(patch function (body))` allows the AI to surgically update specific functions without rewriting the whole file. |
| 28 | **Built-in Rate Limiting / Circuit Breakers** | Pending | 4.0 (4×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | Native `(rate_limit "10/s" (fetch ...))` provides essential guardrails against AI DDoS or loops. |
| 29 | **Implicit Context Threading** | Pending | 3.0 (3×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | `(with_context db ...)` auto-generates Go code that threads dependencies implicitly, saving cognitive load. |
| 30 | **String Manipulation Suite (`str_split`, `str_join`, `regex`)** | Pending | 7.0 (7×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | Essential for parsing and lexing text, required for self-hosting. |
| 31 | **Mutable Collections (`append`, `map_set`)** | Pending | 8.0 (8×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | Needed to build up dynamic lists (like AST children) and manage state. |
| 32 | **Advanced Control Flow (`while`, `match`)** | Pending | 6.5 (6.5×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | State machines and parsers require `while` loops and pattern matching for tokens. |
| 33 | **Full File System I/O (`write_file`, `mkdir`)** | Pending | 6.0 (6×1.0÷1) | Sonnet 3.5 | Gemini 1.5 Pro | Necessary for the transpiler to write out generated `.go` files and manage projects. |

### 26. LLM-Native Primitives
* **Description:** Add built-in nodes like `(llm_generate "prompt" model="...")` and `(vector_embed text)`.
* **Why:** Makes it trivial for an AI to write applications that spawn or utilize other AIs.
* **Impact:** 9/10 (Critical for an AI-first language).

### 27. AST-Level Semantic Patching
* **Description:** Introduce a `(patch function_name (new_body))` node.
* **Why:** LLMs struggle with rewriting large files perfectly. A patch node would allow surgical updates.
* **Impact:** 8/10 (High).

### 28. Built-in Rate Limiting / Circuit Breakers
* **Description:** Add a native `(rate_limit "10/s" (fetch ...))` or `(retry 3 (fetch ...))`.
* **Why:** AI agents writing automation can accidentally DDoS APIs or fall into infinite loops.
* **Impact:** 7/10 (High).

### 29. Implicit Context Threading
* **Description:** A `(with_context db ...)` block that auto-generates Go code threading dependencies.
* **Why:** Removes the need for the AI to remember to pass `req`, `db`, or `context.Context` to every sub-function.
* **Impact:** 6/10 (Medium).

### 30. String Manipulation Suite
* **Description:** Add standard string operations such as `(str_split s sep)`, `(str_join list sep)`, `(str_sub s start end)`, and `(regex_match pattern s)`.
* **Why:** Self-hosting a transpiler requires reading and manipulating raw text efficiently (e.g., tokenizing source code).
* **Impact:** 8/10 (High - blocking for self-hosting).

### 31. Mutable Collections
* **Description:** Add `(append list item)`, `(map_set dict key val)`, and `(map_delete dict key)` to mutate data structures after creation.
* **Why:** The AST builder needs to push parsed child nodes into an array dynamically. Currently, only static lists exist.
* **Impact:** 9/10 (Critical - blocking for self-hosting).

### 32. Advanced Control Flow
* **Description:** Introduce `(while cond body)` for unbounded loops, and `(match var (val body)...)` for cleanly branching on token types.
* **Why:** Writing state machines (Lexers and Parsers) with just basic `for` range loops is extremely difficult.
* **Impact:** 7/10 (High - blocking for self-hosting).

### 33. Full File System I/O
* **Description:** Expand the planned file I/O to include robust writing and directory management: `(write_file path data)`, `(mkdir path)`, and `(list_dir path)`.
* **Why:** A compiler needs to manage projects, traverse source directories, and write output binary/code files to disk.
* **Impact:** 8/10 (High - blocking for self-hosting).
