# ssg

A document transformation pipeline and composable static site generator toolkit for Go.

`ssg` provides the building blocks for a content pipeline: load pages from files
(or any source), enrich them with cross-page data (navigation, taxonomies), and
render them through a chain of transformations to produce HTML output or whatever you want.

```
[source]  →  []ContentSourceConfig  →  Render(pipeline, pages, globals)
```

## Why? And Alternatives

There are many static site generators!

The most focus on a "no programming required" model which limits extensibility

The closest model is [metalsmith](https://metalsmith.io)

## Core concepts

**`Renderer`** — a single transformation step:

```go
type Renderer func(wr io.Writer, src io.Reader, data any) error
```

**`ContentSourceConfig`** — page metadata as a plain map with typed accessors
(`.OutputFile()`, `.TemplateName()`, `.InputFile()`, `.Get(key)`).

**Pipeline** — `Render` chains renderers using buffer-swapping. Each step reads
the previous output and writes to a fresh buffer. No goroutines.

## Usage

```go
loadConf := ssg.LoadConfig{
    ContentDir:      "content",
    InputExt:        ".md",
    PathTransformer: ssg.CleanURLs(".md", ".html"),
    MetaSplit:       ssg.MetaSplitYaml,
    MetaParser:      ssg.MetaParseJson,
}

pipeline := []ssg.Renderer{
    ssg.Must(ssg.NewPageRender("layout", fns)),
    ssg.WriteOutput("public"),
}

var pages []ssg.ContentSourceConfig
ssg.LoadContent(loadConf, &pages)

ssg.Render(pipeline, pages, map[string]any{
    "Site": map[string]any{"BaseURL": "https://example.com"},
})
```

### Taxonomy pages

```go
byTag := ssg.GroupByStrings(pages, "Tags")
for tag, tagPages := range byTag {
    pages = append(pages, ssg.NewPage(
        "tags/"+slug(tag)+"/index.html", "tag-list/index.html",
        map[string]any{"Tag": tag, "Pages": tagPages},
    ))
}
```

### Path transformers

| Function | Example |
|---|---|
| `CleanURLs(".md", ".html")` | `posts/foo.md` → `posts/foo/index.html` |
| `UglyURLs(".md", ".html")` | `posts/foo.md` → `posts/foo.html` |
| `SlugNormalize(next)` | lowercases and hyphenates before applying next |

## Sub-modules

Each sub-module is a separate Go module and can be imported independently.

### Renderers (`render/`)

| Module | Import path | Description |
|---|---|---|
| **htmlclean** | `github.com/client9/ssg/render/htmlclean` | Normalizes HTML fragments via `golang.org/x/net/html`; prevents malformed content from breaking page layout |
| **markdown** | `github.com/client9/ssg/render/markdown` | Markdown → HTML via Goldmark with GFM and auto heading IDs; `markdown.New()` or `markdown.NewGoldmark(g)` |
| **minify** | `github.com/client9/ssg/render/minify` | Minifies HTML/CSS/JS/SVG via `tdewolff/minify`; MIME type derived from output file extension |
| **shortcode** | `github.com/client9/ssg/render/shortcode` | Embedded `$cmd[args]{body}` macro engine; register handlers by name, errors accumulate without stopping rendering |

The shortcode syntax:

```
$cmd
$cmd[arg1 arg2]
$cmd[name=value key="val"]
$cmd{body}
$cmd[args]{body}
$$   →  literal $
```

### Metadata parsers (`meta/`)

| Module | Import path | Description |
|---|---|---|
| **email** | `github.com/client9/ssg/meta/email` | Email-style `Key: Value` headers with optional type coercion |
| **yaml** | `github.com/client9/ssg/meta/yaml` | YAML frontmatter via `gopkg.in/yaml.v3` |
| **toml** | `github.com/client9/ssg/meta/toml` | TOML frontmatter via `github.com/BurntSushi/toml` |

Splitters built into the root module: `MetaSplitYaml`, `MetaSplitJson`, `MetaSplitToml`, `MetaSplitEmail`.

### Template functions (`tmpl/`)

| Module | Import path | Description |
|---|---|---|
| **stdfuncs** | `github.com/client9/ssg/tmpl/stdfuncs` | Stdlib-only `template.FuncMap`; no external dependencies. Covers strings, math, collections, path, time, encoding, safe types, and URL helpers |

```go
t := template.New("page").Funcs(stdfuncs.FuncMap())

// Combine with your own:
fns := stdfuncs.Merge(stdfuncs.FuncMap(), template.FuncMap{"myFunc": myFunc})
```

### Tools (`cmd/`)

| Command | Description |
|---|---|
| `cmd/swapfrontmatter` | Convert frontmatter between YAML, JSON, and email formats; flags: `-from`, `-to`, `-write` |

## Sample

`sample/` is a complete working site: JSON frontmatter, HTML content with
`text/template` macros, page templates, tag taxonomy, and HTML pretty-printing.

```bash
cd sample && make run   # renders to sample/public/
```

## Development

```bash
make test    # go test ./...
make lint    # go mod tidy, gofmt, golangci-lint
make env     # install golangci-lint, goimports
```

Sub-modules each have their own `go.mod`. Run `go test ./...` from their directory,
or use the workspace: `go work sync` at the repo root.
