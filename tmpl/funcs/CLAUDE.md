# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go test ./...          # run all tests
go test -run TestName  # run a single test
```

## Purpose

A stdlib-only `template.FuncMap` for Go's `text/template` and `html/template`. No external dependencies — ever. If a proposed function requires a third-party import, it does not belong here.

## API shape

Two exported symbols:

- `FuncMap() template.FuncMap` — returns all functions in one map, ready to pass to `template.New().Funcs()`
- `Merge(fms ...template.FuncMap) template.FuncMap` — combines maps; later entries win on collision

## Argument order convention

**Subject first, matching Go stdlib.** `strings.Contains(s, substr)` → `contains $s "sub"` in templates. Do not use pipeline-optimized order (subject last). Single-argument functions (`upper`, `lower`, `trim`) work in pipelines regardless.

## Adding functions

New functions go in the appropriate file (`strings.go`, `math.go`) or a new file if a new category is introduced. Register them in the file's unexported `*FuncMap()` helper, not directly in `FuncMap()`.

### Math functions

All math functions accept `any` and return `(float64, error)`. Use `toFloat64` for conversion — it handles all integer widths, both float sizes, and numeric strings. Do not add integer-returning variants; callers use `printf` for formatting.

### String functions

Prefer direct assignment of stdlib functions (`strings.ToLower`, etc.) over wrappers. Only write a custom function when stdlib has no direct equivalent (e.g. `truncate`).

## Documentation

Individual functions registered in a `template.FuncMap` are typed as `any` and do
not appear in `go doc` output. The convention here is **named exported wrappers**:
write an exported function (e.g. `Truncate`, `ReplaceAll`) that is both registered
in the map and callable directly from Go. This surfaces the signature on pkg.go.dev
and makes the function usable without a FuncMap.

Pair each exported function with an `Example` test in `example_test.go` — these
render on pkg.go.dev and serve as runnable documentation.

See TODO.md for the backlog of planned functions and the full documentation plan.

## Naming

The module may be renamed from `tmpl/funcs` to `tmpl/stdfuncs` if additional funcmap varieties (e.g. ones with third-party deps) are added. Raise this before adding any new module under `tmpl/`.
