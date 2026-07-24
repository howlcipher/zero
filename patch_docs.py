import sys

# Patch improvements.md
target_file = "improvements.md"
with open(target_file, "r") as f:
    content = f.read()

content = content.replace("| 46 | [Close the Benchmark-Found Gaps and Make It a Standing Metric](#46-close-the-benchmark-found-gaps-and-make-it-a-standing-metric) | Pending |",
                          "| 46 | [Close the Benchmark-Found Gaps and Make It a Standing Metric](#46-close-the-benchmark-found-gaps-and-make-it-a-standing-metric) | Done |")

old_text = "8. **Autonomous `agy --mode accept-edits` calls can be blocked by the Claude Code permission classifier.** In auto-mode sessions, invoking `agy -p \"...\" --mode accept-edits ...` from Bash can be denied outright by the session's auto-mode classifier (observed 2026-07-23), even though the same command works when the user is prompted interactively. When this happens, do not retry the identical call — either fall back to `--mode manual`/a mode that surfaces edits for review, ask the user to approve the Bash permission rule, or, for genuinely trivial and fully-specified diffs (e.g. a one-clause fix with an exact fix sketch already in the backlog), apply the edit directly with Edit/Write instead of delegating."
new_text = old_text + "\n9. **Benchmark Regression Gate.** Any transpiler change touching `defun`/`type_hint`, `read_file`, `str_split`, or the `test` block must re-run the benchmark harness in `benchmarks/language_write_cost/`, update `results.csv` and `docs/language_write_cost_benchmark.md`'s tables, and note the delta."

content = content.replace(old_text, new_text)

with open(target_file, "w") as f:
    f.write(content)

# Patch bugs.md
target_file2 = "bugs.md"
with open(target_file2, "r") as f:
    content2 = f.read()

content2 = content2.replace("| 17 | [No string-to-number parsing primitive](#17-no-string-to-number-parsing-primitive) | Done (2026-07-23) |", "| 17 | [No string-to-number parsing primitive](#17-no-string-to-number-parsing-primitive) | Done (2026-07-23) |")

with open(target_file2, "w") as f:
    f.write(content2)

print("Markdown updated!")
