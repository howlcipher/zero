# Zero Transpiler

Zero is an AI-first, Lisp-like coding language designed specifically to be written by Large Language Models (LLMs). It transpiles directly into robust, production-ready Go.

## Why Zero? (Why not just write Go?)

You might wonder: *If Zero just compiles into Go, why not have the AI write Go directly?*

1. **Hallucination-Proof Generation**: Modern LLMs often hallucinate invalid syntax or complex abstractions in strictly typed languages. Zero uses simple, uniformly structured S-expressions (Lisp-like grammar). Because the syntax is so simple, we can use tools like `outlines` (as seen in `orchestrator.py`) to mathematically guarantee the AI generates perfectly balanced, structurally valid code.
2. **Immediate Semantic Feedback**: If the AI attempts to do something semantically invalid (e.g., calling a method that doesn't exist), the Go transpiler immediately catches it and returns a clean, localized JSON error. The Orchestrator automatically feeds this back to the AI for self-correction.
3. **Abstraction Constraints**: The AI is strictly constrained by what the transpiler supports. It cannot hallucinate complex, dangerous, or unintended behavior unless an explicit AST mapping exists for it.

Zero combines the **predictability and simplicity of S-expressions** (for the AI to write) with the **performance, safety, and ecosystem of Go** (for the server to run).

## Project Roadmap & Future Home

Zero has evolved from a local script into a standalone transpiler toolchain. To truly achieve its goal of being an AI-first language, **Zero is slated to be moved into its own independent repository**. 

This will allow:
- Seamless installation by AI agents anywhere via `curl` and GitHub Releases.
- Proper semantic versioning and cross-platform CI/CD pipelines.
- A dedicated standard library (`std.zero`, `http.zero`).
- Focus on pure AI optimizations like LLM-native primitives, AST-level semantic patching, and implicit context threading without being tied to a general knowledge library.

## Installation & Requirements

To write and run Zero manually, you only need **Go**. To use the AI Orchestrator to generate Zero code, you need a few additional tools.

### Prerequisites
1. **Go 1.20+**: Required to compile the transpiled Go code.
2. **Python 3.10+**: Required for the AI Orchestrator script.
3. **Ollama**: Required to run local LLMs (like `llama3`) for generating code.

### Setup Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/howlcipher/zero.git
   cd zero
   ```
2. (Optional but recommended) Set up a Python virtual environment:
   ```bash
   python -m venv venv
   source venv/bin/activate
   ```
3. Install the required Python packages for the orchestrator:
   ```bash
   pip install outlines openai
   ```
4. Start Ollama and download the Llama 3 model (in a separate terminal):
   ```bash
   ollama serve
   ollama pull llama3
   ```

## Hello World Example

Here is a basic HTTP server written in Zero that serves a text response and a JSON endpoint.

Create a file called `hello.zero`:

```lisp
(http_server 8080
  (route "/" (lambda (req)
    (res 200 "text/plain" "Hello, World! Zero language is alive!")
  ))
  
  (route "/json" (lambda (req)
    (let (msg (dict ("status" "success") ("message" "Hello from Zero JSON endpoint!")))
      (res_json 200 msg)
    )
  ))
)
```

### CLI "Hello World" Example

If you want to build a simple command-line script instead of a web server, you can use the `(cli_app ...)` root block instead:

Create `cli_hello.zero`:

```lisp
(cli_app
  (print "Hello, World!")
  (let (name "Zero")
    (print "Welcome to" name)
  )
)
```

### Command Line Arguments

Zero scripts can effortlessly read command-line parameters using `(cli_args)`:

```lisp
(cli_app
  (print "All arguments:" (cli_args))
  (print "First argument:" (cli_args 0))
)
```

### AI Orchestration Example

Zero comes with built-in primitives to orchestrate other AIs natively and enforce constraints effortlessly:

```lisp
(cli_app
  (try_let (resp (llm_generate "Translate 'Hello World' to French" "llama3"))
    (catch err (print "Error:" err))
    (print "AI says:" resp)
  )
  
  (struct User (name string) (age int))
  
  ;; Coerce messy text into a strict struct
  (try_let (user_struct (fuzzy_cast User "{ \"name\": \"Alice\", \"age\": 30 }"))
    (catch err (print err))
    (print user_struct)
  )
  
  ;; Intent-based qualitative validation
  (let (is_valid (assert_semantic "Alice is a doctor" "is professional"))
    (if (= is_valid true)
      (print "Approved")
      (print "Rejected")
    )
  )
)
```

### AST-Level Semantic Patching

Zero supports surgically updating functions without rewriting the entire file, which is highly beneficial for LLMs struggling with large file generation:

```lisp
(cli_app
  (defun foo (x) (return "Old behavior: "))
  (patch foo (return "New patched behavior: "))
  (let (v (call foo "test")) (print v))
)
```

### Implicit Context Threading

Zero can automatically inject context variables into function calls within a specific block, reducing the cognitive load for AIs to remember to thread variables like `req`, `db`, or `ctx`:

```lisp
(cli_app
  (defun fetch_user (db user_id)
    (print "Fetching user" user_id "from" db)
  )
  (let (db "PostgreSQL")
    (with_context (db)
      ;; Automatically expanded to (call fetch_user db 123)
      (call fetch_user 123)
    )
  )
)
```

## How to Run

1. **Transpile and Run in one step**:
   To immediately transpile your `.zero` code into Go and execute it, run:
   ```bash
   go run zero.go hello.zero && go run server.go
   ```

2. **Build a Standalone Binary**:
   If you want to compile the transpiled Go code into a highly optimized, standalone binary:
   ```bash
   # 1. Transpile to server.go
   go run zero.go hello.zero
   
   # 2. Compile Go into an executable
   go build -o hello server.go
   
   # 3. Run the binary
   ./hello
   ```

The server will spin up on `http://localhost:8080`.

## Generating Code with AI (Orchestrator)

Zero is designed to be written by an AI. We provide an orchestrator script (`orchestrator.py`) that handles the interaction with the LLM, strictly enforces syntax boundaries using `outlines`, and handles error feedback loops.

1. Ensure Ollama is running (`ollama serve`) and the `llama3` model is available.
2. Open `orchestrator.py` and modify the `prompt` variable to instruct the AI on what to build.
   ```python
   prompt = "Build a web server on port 8080 with a root route returning 'root' and an /api route returning 'api'."
   ```
3. Run the orchestrator:
   ```bash
   python orchestrator.py
   ```
4. The AI will generate a `.zero` file (by default `app.zero`). The orchestrator will automatically run the Go transpiler.
5. **Self-Correction loop**: If the transpiler encounters a semantic error (e.g. invalid arguments or missing variables), it outputs a localized JSON error. The orchestrator intercepts this error and sends it back to the AI for automatic self-correction.
6. Once transpilation succeeds, the orchestrator compiles the Go binary and executes the newly generated application.

### Automation and Advanced Control Flow

Zero has native support for file operations, subprocess execution, advanced loops (`while`, `match`), and string manipulation for easy automation scripting and state-machine building:

```lisp
(cli_app
  (write_file "hello.txt" "Hello from Zero!")
  (try_let (content (read_file "hello.txt"))
    (catch err (print err))
    (print content)
  )
  (exec "rm" "hello.txt")
)
```

### Native Unit Test Blocks

Zero supports Test-Driven Development natively. You can include `(test "description" ...)` blocks in your code, which the transpiler will extract and convert directly into Go test functions (`_test.go`). This allows AIs to iterate rapidly with test-driven workflows:

```lisp
(cli_app
  (defun add (a b)
    (return (+ a b))
  )
  
  (test "add function returns correct sum"
    (let (result (call add 2 3))
      (if (!= result 5)
        (print "Error: expected 5 got" result)
      )
    )
  )
)
```

