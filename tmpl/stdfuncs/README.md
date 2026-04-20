# stdfuncs - additional functions for Go templates

The default functions in `text/template` and `html/template` are minimal. This extends them.

## What's Included?

- strings — case, trim, search, replace, split/join, truncate, rune-aware length
- math — arithmetic, rounding, min/max, clamp, pow
- cast — `toInt`, `toFloat` for frontmatter values that arrive as strings
- encoding — `jsonify` for embedding data in `<script>` blocks
- date and time — `now`, `parseTime`; use `time.Time` methods for formatting
- url / safe types — `urlEncode`, `urlPathEscape`, `safeHTML`, `safeCSS`, etc.
- path — `pathBase`, `pathDir`, `pathExt`, `pathJoin`, `pathClean`
- lists (slices) — `list`, `seq`, `take`, `drop`, `sort`, `sortNum`, `reverse`, `concat`, …
- dicts (maps) — `dict`, `keys`, `values`, `merge`

## Goals

- **Independent and exportable** — serves as a base, or for use in different templating systems
- **Stdlib only** — keep it simple; functions requiring external deps go in a different module
- **Not pipeline-based** — pipeline order looks elegant for single-argument functions, then gets confusing. Argument order follows Go stdlib (subject first).
- **Minimal** — want a descending sort? Use `reverse (sort $list)`. Want all but last? Use `drop $list -1`.
- **Prefer separate functions over extra arguments** — `sort` and `sortNum` instead of a mode flag
- **Immutable data structures** — all functions return new values, never modify inputs

## Alternatives

[masterminds/sprig](https://github.com/masterminds/sprig) — appears semi-abandoned, pipeline-based, has a number of unusual functions.

[Hugo](https://gohugo.io/) — the static site generator has many functions, but inconsistent design and argument order optimized for pipelines. Implementation is tightly coupled to Hugo internals.

## Not Included

- **Internationalization / titlecase** — requires `golang.org/x/text`; good for a separate module. Note: all string operations here are rune-aware.
- **Regular expressions** — defer until use cases in templates are better understood
- **Base64 encoding** — two competing encodings (standard vs URL-safe); add when use case is clear
- **Random / shuffle** — non-deterministic output is hostile to SSG reproducibility
- **Checksum and hashes** — limited uses, many variations; good for a separate module
- **Cryptography** — limited use, many variations
- **OS and environment** — pass these as data to the template instead
- **Math trig** — limited utility in HTML templates

## Recipes

**Descending sort**
```
reverse (sort $list)
```

**Capitalize first letter only** (leave rest unchanged)
```
firstUpper $str
```

**Capitalize first letter, lowercase the rest** (like Jinja2 `capitalize`)
```
capitalize $str
```

**Last N elements**
```
take $list -2
```

**All but last N elements**
```
drop $list -2
```

**Join with special last separator** ("a, b and c")
```
printf "%s and %s" (join (drop $list -1) ", ") (last $list)
```

**Add trailing separator**
```
printf "%s," (join $list ", ")
```

**Chomp** (remove trailing newline)
```
trimRight $str "\r\n"
```

**toString**
```
printf "%v" $val
```

**Date formatting**
```
{{ $date.Format "2006-01-02" }}
```

**UTC time**
```
{{ now.UTC }}
```

**Simple loop counter** (Go 1.24+)
```
{{ range 10 }}{{ . }}{{ end }}
```
outputs 0–9. Use `seq` for non-zero start or step: `{{ range (seq 3 7) }}`.

**String Pad Right**
```
{{ printf "%20s" $atr }}
```

**String Pad Left**
```
{{ printf "%-20s" $astr }}
```


