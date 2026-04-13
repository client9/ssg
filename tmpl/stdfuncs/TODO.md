# TODO

## Remaining work

_Nothing pending._

---

## Completed

| Category | File | Functions |
|---|---|---|
| Strings | `strings.go` | `lower`, `upper`, `trim`, `trimPrefix`, `trimSuffix`, `trimLeft`, `trimRight`, `contains`, `hasPrefix`, `hasSuffix`, `count`, `replace`, `replaceAll`, `repeat`, `split`, `join`, `fields`, `lenRunes`, `truncate`, `firstUpper`, `capitalize` |
| Math | `math.go` | `add`, `sub`, `mul`, `div`, `mod`, `abs`, `ceil`, `floor`, `round`, `min`, `max`, `pow`, `modBool`, `clamp` |
| Path | `path.go` | `pathBase`, `pathDir`, `pathExt`, `pathJoin`, `pathClean` |
| Safe / URL | `safe.go` | `safeCSS`, `safeHTML`, `safeHTMLAttr`, `safeJS`, `safeJSStr`, `safeURL`, `urlEncode`, `urlPathEscape` |
| Encoding | `encoding.go` | `jsonify` |
| Cast | `cast.go` | `toInt`, `toFloat` |
| Time | `time.go` | `now`, `parseTime` |
| Collections | `collections.go` | `list`, `dict`, `seq`, `first`, `last`, `take`, `drop`, `reverse`, `compact`, `concat`, `sort`, `sortNum`, `where`, `keys`, `values`, `merge`, `in`, `default`, `cond` |
| Documentation | `example_test.go` | Example tests for all exported functions |

---

## Design notes

### Collections

- `list` not `slice` — avoids shadowing the built-in `slice` (subslicing) action
- `compact` = consecutive-duplicate removal only (`slices.Compact` semantics); full dedup = `compact (sort $list)`
- `sort`/`sortNum` are separate; use `reverse` for descending — no direction flag
- `sort` is type-aware: numeric types sort numerically, `time.Time` sorts chronologically, everything else lexicographically; for `[]any` the first non-nil element determines mode
- `where` is equality-only; comparison operators become new functions (`whereGt`, `whereLt`) if needed
- `take`/`drop` accept negative n: `take $list -2` = last 2, `drop $list -2` = all but last 2
- `seq` is 1-based: `seq 5` → `[1 2 3 4 5]`; `seq 3 7` → `[3..7]`; `seq 1 10 2` → `[1 3 5 7 9]`; `{{ range n }}` (Go 1.24) covers the simple counter case
- `merge` template key maps to exported `MergeMaps` to avoid collision with `funcs.Merge` (FuncMap combiner)
- Collection functions return new structures — never mutate inputs

### Strings

- `capitalize` = `firstUpper(lower(s))` — opt-in destructive lowercasing, unlike `firstUpper` which leaves other chars unchanged
- `replace` defaults to n=1; `replace $s $old $new -1` = `replaceAll`
- `lenRunes` counts runes not bytes — consistent with `truncate`, `take`, `drop` which are all rune-aware

### Math

- `clamp $val $min $max` — subject first (value being clamped), bounds follow; errors on min > max

---

## Not adding

| Function | Reason |
|---|---|
| `strings.Title` | Deprecated in Go; correct titlecase requires `golang.org/x/text` (no external deps) |
| `replaceRE` / `findRE` / `findRESubmatch` | Defer until regexp use cases in templates are better understood |
| `CountWords` / `CountRunes` | Low priority; add if reading-time estimates are needed |
| `math` trig (`sin`, `cos`, `tan`, …) | Almost never used in HTML templates |
| `math.rand` / `shuffle` | Non-deterministic output is hostile to SSG reproducibility |
| `path.Split` | Returns two values; awkward in templates — use `pathDir`/`pathBase` instead |
| `time.In` / `time.ParseDuration` / `time.Duration` | Niche; add when there is a real use case |
| `dateFormat` / `year` / `month` / `day` | Redundant — `time.Time` methods work directly in templates: `{{.Date.Format "2006-01-02"}}`, `{{.Date.Year}}` |
| `nowUtc` | Redundant — `{{now.UTC}}` works |
| `toString` | Redundant — `printf "%v" $x` or `{{ $x }}` covers it |
| `chomp` | Redundant — `trimRight $s "\r\n"` covers it |
| `base64Encode` / `base64Decode` | Two competing encodings (standard vs URL-safe); no clear use case yet |
| `append` / `set` | Workarounds exist (`concat`, `merge`); mutation semantics require returning new structures — add when there is a real use case |
| `apply` | Requires stringly-typed function names — no clean implementation without first-class functions |
| `groupBy` | Better handled before the template sees the data; covered by `ssg.GroupByString`/`GroupByStrings` |
| `delimit` (join with last-sep) | Composable: `concat (join (take $list -1) ", ") " and " (last $list)` |
