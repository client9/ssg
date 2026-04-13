# TODO

## Remaining work

_Nothing pending._

---

## Notes

Collections design notes (decisions not obvious from the code):
- `list` not `slice` — avoids shadowing the built-in `slice` (subslicing) action
- `compact` = consecutive-duplicate removal only (`slices.Compact` semantics); full dedup = `compact (sort $list)`
- `sort`/`sortNum` are separate; use `reverse` for descending — no direction flag
- `where` is equality-only; if comparison operators are ever needed they become new functions (`whereGt`, `whereLt`, etc.) following the same pattern
- `take`/`drop` are count-based, collection-first: `take $list 5`, `drop $list 3`
- `seq` is 1-based: `seq 5` → `[1 2 3 4 5]`; `seq 3 7` → `[3..7]`; `seq 1 10 2` → `[1 3 5 7 9]`
- `merge` template key maps to exported `MergeMaps` to avoid collision with `funcs.Merge` (FuncMap combiner)
- `sort` does not handle `time.Time`; ISO 8601 strings (`"2006-01-02"`) sort correctly lexicographically

---

## Not adding

| Function | Reason |
|---|---|
| `strings.Title` | Deprecated in Go; correct titlecase requires `golang.org/x/text` (no external deps) |
| `strings.FindRE` / `FindRESubmatch` | Defer until there is a concrete use case |
| `CountWords` | Low priority; add if reading-time estimates are needed |
| `math` trig (`sin`, `cos`, `tan`, …) | Almost never used in HTML templates |
| `path.Split` | Returns two values; awkward in templates — use `pathDir`/`pathBase` instead |
| `time.In` / `time.ParseDuration` / `time.Duration` | Niche; add when there is a real use case |
| `dateFormat` / `year` / `month` / `day` | Redundant — `time.Time` methods work directly in templates: `{{.Date.Format "2006-01-02"}}`, `{{.Date.Year}}` |
| `toString` | Redundant — `printf "%v" $x` or just `{{ $x }}` covers it |
| `chomp` | Redundant — `trimRight $s "\r\n"` covers it |
| `replaceRE` / `findRE` / `findRESubmatch` | Defer until regexp use cases in templates are better understood |
| `base64Encode` / `base64Decode` | Two competing encodings (standard vs URL-safe); no clear use case yet to pick one |
