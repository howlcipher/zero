# Cross-Language "AI Write Cost" Benchmark

Traditional language benchmarks measure compile time and runtime speed. Zero doesn't compete on those — it transpiles to Go and inherits Go's runtime performance, unmodified. Zero's actual pitch (see "Why Zero?" in the [README](../README.md)) is that its constrained, uniform S-expression grammar is cheaper for an LLM to *write correctly* than a full-size language — fewer hallucinated syntax errors, fewer self-correction round-trips, less output.

This benchmark measures that claim directly: for a fixed set of task prompts, how long (wall-clock) and how many tokens does it take an LLM to produce a **correct, compiler/runtime-verified** solution in Zero vs. Go, Python, Node.js, C#, and Java?

This is improvement [#44](../improvements.md#44-add-cross-language-ai-write-cost-benchmark) in the project backlog.

## Methodology

- **3 fixed tasks**, implemented once per language (18 programs total):
  - **A — Hello World HTTP+JSON server**: root route returns plain text, `/json` returns a JSON body. Mirrors the README's existing Zero example.
  - **B — CLI file-parsing tool**: read a file of names (one per line), print a greeting per non-blank line, handle a missing file gracefully (print an error, don't crash/stack-trace).
  - **C — Function + unit test**: an `add(a, b)` function plus an idiomatic unit test asserting `add(2, 3) == 5`.
- **Write time**: wall-clock seconds from immediately before drafting a solution to immediately after the source file(s) were written, timestamped with `date +%s.%N` bracketing each generation. This measures the full cost of producing the code — including this session's reasoning and tool-call overhead — not raw model decode latency. It is **not** a controlled lab benchmark; treat it as a reproducible proxy, not a precision timer.
- **Token count**: the final source's token count under `tiktoken`'s `cl100k_base` encoding, as a reproducible proxy for LLM output-token cost. Project-scaffolding files that a real developer would generate via `dotnet new`/`go mod init`/etc. rather than hand-write (`.csproj`, `go.mod`, `package.json`) are excluded from the count; hand-written source and test files are included.
- **Verification**: every single one of the 18 programs was actually compiled and run (not just reviewed) before its numbers were recorded — `go build`/`go run` (Zero, Go), `python3`/`pytest` (Python), `node`/`node --test` (Node.js), `dotnet build`/`dotnet test` (C#), `javac`/JUnit console (Java). Java (`openjdk`) and .NET (`dotnet`) were not installed on the benchmark machine beforehand; both were installed via Homebrew (user-space, no `sudo`/`rpm-ostree` needed on this Bazzite/Kinoite host) specifically so all 6 languages got equal, compiler-verified footing.
- Raw timestamps, per-task/per-language results, and all 18 source programs are in [`benchmarks/language_write_cost/`](../benchmarks/language_write_cost/) (`results.csv` has the raw data this doc summarizes).

## Results

### A — Hello World HTTP+JSON

| Language | Write time (s) | Tokens | Verified |
|---|---|---|---|
| Zero | 5.3 | 84 | pass |
| Go | 5.1 | 144 | pass |
| Python | 4.6 | 164 | pass |
| Node.js | 3.9 | 111 | pass |
| C# | 5.4 | 177 | pass |
| Java | 6.9 | 237 | pass |

Zero's smallest and clearest win: this is its actual designed niche (HTTP/JSON web handlers), and it produces the fewest tokens of any language by a wide margin while writing just as fast as the others.

### B — CLI file-parsing tool

| Language | Write time (s) | Tokens | Verified |
|---|---|---|---|
| Zero | 39.9 | 88 | pass |
| Go | 8.6 | 97 | pass |
| Python | 4.0 | 61 | pass |
| Node.js | 4.6 | 87 | pass |
| C# | 4.4 | 80 | pass |
| Java | 6.8 | 127 | pass |

Zero's token count is still mid-pack-good here, but its write time is an outlier — **not** because the final Zero program is large or complex (it's the second-shortest of the six), but because writing it surfaced two real, previously-undocumented transpiler gaps mid-task:

- Zero has no deterministic string-to-number parsing primitive (filed as [bug #17](../bugs.md#17-no-string-to-number-parsing-primitive)). The original Task B design ("sum a numeric column from a file") turned out to be unwritable in Zero without abusing the LLM-backed `fuzzy_cast`, so Task B was redesigned around string-only processing to keep the comparison fair.
- Even the string-only version didn't compile on the first attempt: `read_file` returns a Go `[]byte`, and there is no primitive to convert it to `string` before passing it to `str_split` (documented as a related finding on bug #17). The fix — `(call string content)` — only works because `call` blindly emits `funcName(args)` for any symbol name; it's an undocumented escape hatch, not an intentional feature.

That investigation-and-fix time is real and reproducible with today's Zero, but it's a one-time discovery cost tied to two specific, now-filed bugs — not an inherent property of the language design. Both are legitimate weaknesses this benchmark exists to surface, not hide.

### C — Function + unit test

| Language | Write time (s) | Tokens | Verified |
|---|---|---|---|
| Zero | 4.9 | 123 | pass |
| Go | 5.7 | 73 | pass |
| Python | 11.6 | 37 | pass |
| Node.js | 8.4 | 72 | pass |
| C# | 13.5 | 63 | pass |
| Java | 7.9 | 74 | pass |

Zero's weakest showing: it's the **most** token-heavy of all six languages here, not the least. Its native `(test ...)` block is fast to write, but the `defun` requires three separate `(type_hint ...)` statements (for `a`, `b`, and the return value) to get a typed Go function — overhead that Python's fully-dynamic `def add(a, b): return a + b` simply doesn't pay, and that even Go only pays once via its normal parameter-list syntax rather than three extra S-expression forms.

### Totals (all 3 tasks)

| Language | Total write time (s) | Total tokens |
|---|---|---|
| Zero | 50.1 | 295 |
| Go | 19.4 | 314 |
| Python | 20.1 | 262 |
| Node.js | 17.0 | 270 |
| C# | 23.2 | 320 |
| Java | 21.7 | 438 |

## Takeaways

- **Zero wins clearly on its home turf** (HTTP/JSON handlers, Task A) and is competitive-to-mid-pack on tokens overall (3rd of 6, behind Python and Node.js) — but it is **not** a uniform win across every kind of program, and this benchmark's job is to show that honestly rather than only report the flattering half.
- **The retry/discovery cost this benchmark set out to measure did show up — just not in the direction assumed.** The pitch was that Zero's constrained grammar reduces retries versus mainstream languages. In practice, at this task size, a capable LLM writes correct Go/Python/Node/C#/Java on the first try just as often as Zero — the extra cost instead came from Zero's smaller, less-documented primitive surface (missing type coercion) rather than from syntax hallucination.
- **Two real bugs got filed as a direct result** ([#17](../bugs.md#17-no-string-to-number-parsing-primitive) and its `read_file`/`[]byte` addendum) — this benchmark is as much a bug-finding exercise as a marketing one, and both should be treated as backlog items to close before re-running this benchmark for a future comparison.
- **Task C's type-hint overhead is a real, measurable tradeoff**, not a one-off: three-line boilerplate per typed function is inherent to how `type_hint` currently works, and would keep showing up in any task involving several typed functions.

*Last run: 2026-07-23. Re-run this benchmark after any transpiler change that touches `defun`/`type_hint`, `read_file`, `str_split`, or the `test` block, since all four are exercised directly above.*
