# Improvement 46: Close the Benchmark-Found Gaps and Make It a Standing Metric

## Actions Taken
1. Fixed `defun` to support `type_hints` config nodes and inline typed arguments (e.g. `(defun add ((a int) (b int)) int ...)`). This significantly reduces boilerplate compared to individual `type_hint` nodes.
2. Verified that `to_int`, `to_float`, `to_string`, and `bytes_to_string` are properly implemented. Addressed a bug where `to_int` and `to_float` returned two values (`(value, error)`) by wrapping them in a closure that ignores the error and returns a single zero-value on failure, allowing seamless use inside `let`.
3. Fixed `to_string` codegen which was incorrectly calling `string(%s)` for everything, which returns ASCII chars for integers. Separated `bytes_to_string` (`string(%s)`) and `to_string` (`fmt.Sprint(%s)`).
4. Added `tests/test_improvement_46.zero` to verify the new syntax and functions. Tests compiled and ran successfully.
5. Updated the Working Protocol in `improvements.md` to formalize the benchmark as a regression gate.
6. Marked Improvement 46 as Done in `improvements.md`. Bug 17 was already marked Done but the implementation was refined.
