# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go test ./...          # run all tests
go test -run TestName  # run a single test
```

## Purpose

A stdlib-only `template.FuncMap` for Go's `text/template` and `html/template`. No external
dependencies — ever. If a proposed function requires a third-party import, it does not belong here.

## API shape

Two top-level exported symbols for callers:

- `FuncMap() template.FuncMap` — returns all functions in one map, ready to pass to `template.New().Funcs()`
- `Merge(fms ...template.FuncMap) template.FuncMap` — combines FuncMaps; later entries win on collision

All individual template functions are also exported as named Go functions (e.g. `Truncate`,
`Sort`, `SafeHTML`) so they are callable directly from Go and visible in `go doc`. See the
Documentation section below.

## File structure

Each category lives in its own file with an unexported `*FuncMap()` helper registered in `funcmap.go`:

| File | Category | Key functions |
|---|---|---|
| `strings.go` | Strings | `lower`, `upper`, `trim*`, `contains`, `split`, `join`, `truncate`, … |
| `math.go` | Math | `add`, `sub`, `mul`, `div`, `mod`, `abs`, `ceil`, `floor`, `round`, `min`, `max`, `pow`, `modBool` |
| `path.go` | Path | `pathBase`, `pathDir`, `pathExt`, `pathJoin`, `pathClean` |
| `safe.go` | Safe types / URL | `safeHTML`, `safeCSS`, `safeJS`, `safeURL`, …, `urlEncode`, `urlPathEscape` |
| `collections.go` | Collections | `list`, `dict`, `seq`, `first`, `last`, `take`, `drop`, `sort`, `sortNum`, `where`, `keys`, `values`, `merge`, `in`, `default`, `cond`, … |

Planned but not yet implemented: `time.go`, `cast.go`, `encoding.go`, and string additions
(`firstUpper`, `chomp`, `replace`, `replaceRE`). See TODO.md.

## Argument order convention

**Subject first, matching Go stdlib.** `strings.Contains(s, substr)` → `contains $s "sub"` in
templates. Do not use pipeline-optimized order (subject last). Single-argument functions
(`upper`, `lower`, `trim`) work in pipelines regardless.

**Exception — `take` and `drop`:** these are collection-first (`take $list 5`, `drop $list 3`)
matching the convention, which also happens to be the natural English reading ("take from list, 5").

## Adding functions

New functions go in the appropriate file (`strings.go`, `math.go`, etc.) or a new file if a new
category is introduced. Register them in the file's unexported `*FuncMap()` helper, not directly
in `FuncMap()`.

### Math functions

All math functions accept `any` and return `(float64, error)` — except `modBool` which returns
`(bool, error)`, `min`/`max` which are variadic `(...any) (float64, error)`, and `pow` which
takes two `any` args. Use `toFloat64` for numeric conversion — it handles all integer widths,
both float sizes, and numeric strings. Do not add integer-returning variants; callers use
`printf` for formatting.

### String functions

Prefer direct assignment of stdlib functions (`strings.ToLower`, etc.) over wrappers. Only
write a custom function when stdlib has no direct equivalent (e.g. `truncate`).

### Collection functions

Use `toSlice(v any) ([]any, error)` (defined in `collections.go`) to accept arbitrary slice
types — it handles `[]any` directly and other slice types via reflection. Use
`reflect.DeepEqual` for equality comparisons across `any` values. Use `slices.SortStableFunc`
for sorting (stable, preserves relative order of equal elements).

The `isZero(v any) bool` helper (in `collections.go`) is the shared definition of "zero value"
used by both `default` and `cond`.

## Documentation

Individual functions registered in a `template.FuncMap` are typed as `any` and do not appear
in `go doc` output. The convention here is **named exported wrappers**: write an exported
function (e.g. `Truncate`, `Sort`) that is both registered in the map and callable directly
from Go. This surfaces the signature on pkg.go.dev and makes the function usable without a FuncMap.

Pair each exported function with an `Example` test in `example_test.go` — these render on
pkg.go.dev and serve as runnable documentation. This is not yet done; see TODO.md.

## Naming

The module may be renamed from `tmpl/funcs` to `tmpl/stdfuncs` if additional funcmap varieties
(e.g. ones with third-party deps) are added. Raise this before adding any new module under `tmpl/`.
