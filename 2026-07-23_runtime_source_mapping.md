# 🐛 Task Journal: Runtime Source Mapping
**Date:** 2026-07-23
**Goal:** Fix Bug #11 (No Runtime Source Mapping)

## Steps:
1.  **Add `Filename` to `Node`:** Update the AST to store the source filename.
2.  **Update Parser:** Update `Parser` and `NewParser` to accept and set the filename.
3.  **Update `expandIncludes` and `main`:** Pass the filename to `NewParser`.
4.  **Inject `//line` directives:** Update `generateCode` and `generateStatement` to prepend `//line <filename>:<line>` to Go statements.
5.  **Test:** Write a test file that triggers a panic, compile it, and verify the panic stack trace maps to the `.zero` file.
