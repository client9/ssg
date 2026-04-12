# TODO

## Planned function groups

### HTML / URL safety
Functions that produce `template.HTML` or `template.URL` typed values, preventing
double-escaping when used with `html/template`.

| Name | Signature | Notes |
|---|---|---|
| `safeHTML` | `(s string) template.HTML` | marks pre-rendered body as trusted |
| `safeURL` | `(s string) template.URL` | marks a URL as trusted |
| `urlEncode` | `(s string) string` | `url.QueryEscape` |
| `urlPathEscape` | `(s string) string` | `url.PathEscape` |

These need a new file `safe.go`. Note that `template.HTML` and `template.URL` are
from `html/template`, not `text/template` — the types are compatible with both.

### Date / time
All stdlib (`time` package). Follow subject-first convention.

| Name | Signature | Notes |
|---|---|---|
| `now` | `() time.Time` | |
| `dateFormat` | `(t time.Time, layout string) string` | Go reference-time layout |
| `year` | `(t time.Time) int` | |
| `month` | `(t time.Time) int` | 1–12 |
| `day` | `(t time.Time) int` | |

New file: `time.go`.

### Collections
Reflection-based helpers for slices. Return `(any, error)` so template execution
stops cleanly on bad input.

| Name | Signature | Notes |
|---|---|---|
| `first` | `(v any, n int) (any, error)` | first n elements of any slice |
| `last` | `(v any, n int) (any, error)` | last n elements |
| `default` | `(def, val any) any` | returns `def` if `val` is the zero value |
| `dict` | `(kvs ...any) (map[string]any, error)` | inline map: `dict "k" v "k2" v2` |

New file: `collections.go`. Use `reflect` — it is stdlib. Decide whether to also
expose a typed `seq(n int) []int` for range loops (no reflection needed).

## Documentation
Every exported function needs a Go doc comment. Currently `FuncMap` and `Merge`
have comments; the individual template functions (registered as `any` in the map)
do not appear in `go doc` output. Options:

- Document each function via a named wrapper (e.g. `func Truncate(s string, n int) string`)
  that is both exported for godoc and registered in the map.
- Add an `example_test.go` with `Example` functions — these render on pkg.go.dev
  and double as runnable tests.
- At minimum, a table in `doc.go` listing every function name, signature, and
  one-line description.

Recommendation: named exported wrappers + `Example` tests. The wrappers also make
the functions directly callable from Go code without going through a FuncMap.
