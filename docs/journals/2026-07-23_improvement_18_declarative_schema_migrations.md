# Improvement 18: Declarative Schema Migrations

**Date:** 2026-07-23
**Status:** Done

## Overview
Implemented declarative schema migrations using a new `(schema ...)` root-level node in `.zero` files.
This node allows developers to define schemas identically to structs (e.g., `(schema "users" (column "id" "int") (column "name" "string"))`), fulfilling the promise of a single source of truth for both data types and SQL DDL.

## Changes Made
1. **Parser & Code Generator Updates (`zero.go`)**:
   - Added support for the `schema` node in `generateCode`.
   - The node parses table names and columns.
   - It outputs a Go `struct` corresponding to the schema (e.g., `Users` for `"users"`).
   - It generates a `CREATE TABLE IF NOT EXISTS` statement and appends it to a global slice `currentSchemaDDLs`.
   
2. **Database Connection Interception (`zero.go`)**:
   - Modified the `db_connect` node processor in `generateStatementRaw`.
   - When a database connection is established, it iterates through all DDL statements in `currentSchemaDDLs` and automatically executes them using `db.Exec`.

3. **Testing**:
   - Wrote `test_schema.zero` inside `tests/` which verifies that the schema definition is executed against a SQLite database and inserts/selects run properly.
   - Verified functionality via `go run zero.go tests/test_schema.zero && go build -o servercheck . && ./servercheck`.

4. **Backlog Updates**:
   - Updated `improvements.md` to mark Improvement 18 as Done.
