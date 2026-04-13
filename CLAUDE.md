# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Test
make test          # go test ./...

# Lint (runs go mod tidy, gofmt -w -s, golangci-lint)
make lint

# Install linting tools
make env           # installs golangci-lint, goimports

# Run a single test
go test -run TestName ./...
```

Sub-packages each have their own `go.mod`. Run `go test ./...` from their directory.

```bash
# Sample site
cd sample && make run    # go run main.go  → outputs to sample/public/
```

## Architecture

This is a **static site generator toolkit** (module `github.com/client9/ssg`). It provides a composable pipeline for transforming content files with frontmatter metadata into HTML output.

### Core types

```go
// A single transformation step in the rendering pipeline
type Renderer func(wr io.Writer, src io.Reader, data any) error

// Page metadata (extracted frontmatter) — a plain map with typed accessors
// (.OutputFile(), .TemplateName(), .InputFile(), .Get(key))
type ContentSourceConfig map[string]any

// LoadConfig controls how content files are found and parsed.
// No role in rendering.
type LoadConfig struct {
    ContentDir      string
    BaseTemplate    string
    MetaSplit       ContentSplitter   // splits raw bytes → (meta, body)
    MetaParser      MetaParser        // parses meta bytes → ContentSourceConfig
    InputExt        string            // file extension filter, e.g. ".md"
    PathTransformer PathTransformer   // maps input path → output path
}

// PathTransformer maps a relative input path to a relative output path.
// Return "" to skip the file. Built-ins: CleanURLs, UglyURLs.
// Compose with SlugNormalize: SlugNormalize(CleanURLs(".md", ".html"))
type PathTransformer func(relPath string) string
```

### Processing model

The only required step is `Render`. Everything before it is user code:

```
[any source]  →  []ContentSourceConfig  →  Render
```

Pages can come from anywhere — files, a database, an API, or constructed directly
with `NewPage`. `LoadContent` is a convenience helper for the common case of loading
from a directory of files; it is not special and can be skipped entirely.

The typical three-stage pattern:

1. **Populate** (`LoadContent` or any source): Build `[]ContentSourceConfig` however
   makes sense. `LoadContent` walks a directory, parses frontmatter, and applies a
   `PathTransformer` to determine output paths.

2. **Enrich** (user code): Compute cross-page data — filter drafts, build navigation,
   generate taxonomy pages with `GroupByStrings` + `NewPage`, sort, group.

3. **Render** (`Render`): Merge `globals` into each page (page frontmatter wins on
   collision), run body through `pipeline` renderers.

```go
loadConf := ssg.LoadConfig{
    ContentDir:      "content",
    InputExt:        ".md",
    PathTransformer: ssg.CleanURLs(".md", ".html"),
}
pipeline := []ssg.Renderer{
    ssg.NewTemplateMacro(fns),
    ssg.Must(ssg.NewPageRender("layout", fns)),
    ssg.WriteOutput("public"),
}

pages := []ssg.ContentSourceConfig{}
ssg.LoadContent(loadConf, &pages)
pages = ssg.FilterPages(pages, isPublished)

// Taxonomy: build tag pages from loaded content
byTag := ssg.GroupByStrings(pages, "Tags")
for tag, tagPages := range byTag {
    pages = append(pages, ssg.SyntheticPage(
        "tags/"+slug(tag)+"/index.html", "tag-list.html",
        map[string]any{"Tag": tag, "Pages": tagPages},
    ))
}

ssg.Render(pipeline, pages, map[string]any{
    "Nav":  buildNav(pages),
    "Site": map[string]any{"BaseURL": "https://example.com"},
})
```

### Pipeline pattern

`MultiRender()` (render.go) chains `Renderer` functions using buffer-swapping — no
goroutines or channels. Each renderer reads from the previous output and writes to a
swapped buffer. Avoid introducing goroutines into the rendering path.

### Path transformation

`paths.go` provides built-in `PathTransformer` implementations:
- `CleanURLs(inputExt, outputExt)` — `foo.md` → `foo/index.html` (default)
- `UglyURLs(inputExt, outputExt)` — `foo.md` → `foo.html`
- `SlugNormalize(next)` — lowercases, replaces spaces/underscores with hyphens, then applies next

`LoadDefaults` sets `CleanURLs(InputExt, ".html")` when `PathTransformer` is nil.

### Taxonomy helpers

`taxonomy.go` provides primitives for generating pages from aggregated metadata:
- `GroupByString(pages, field)` — group by a single-value string field (e.g. `"Category"`)
- `GroupByStrings(pages, field)` — group by a multi-value field (e.g. `"Tags"`);
  handles `[]string`, `[]any`, and bare `string` values
- `SyntheticPage(outputFile, templateName, data)` — create a page not backed by a file;
  sets `Content: []byte{}` so `Render` can process it

### Metadata formats

Splitters (`meta.go`) separate raw file bytes into metadata + body:
- `MetaSplitYaml`, `MetaSplitJson`, `MetaSplitToml`, `MetaSplitEmail`

Parsers convert the metadata bytes into a `ContentSourceConfig`:
- `MetaParseJson` — used for both JSON and YAML (YAML converts to JSON internally)
- `meta/email.Parser(transformers...)` — email-style headers with optional type coercion
- `meta/yaml.Parser()`, `meta/toml.Parser()`

### Sub-packages

Each is a separate Go module. All are included in `go.work` for local development.

**Renderers** (`render/`) — implement `Renderer` and are used as pipeline steps:
- **`render/htmlclean`** — normalizes HTML fragments via `golang.org/x/net/html`
- **`render/markdown`** — Goldmark-based Markdown→HTML with GFM and auto-heading IDs; use `markdown.New()` or `markdown.NewGoldmark(g)`
- **`render/minify`** — minifies HTML/CSS/JS/SVG output via `tdewolff/minify`; MIME type derived from `ContentSourceConfig.OutputFile()` extension using a hardcoded map (not `mime.TypeByExtension` — OS databases are unreliable)

**Metadata parsers** (`meta/`) — implement `MetaParser` for frontmatter:
- **`meta/email`** — email-style `Key: Value` headers; `email.Parser(transformers...)` returns a `MetaParser`
- **`meta/yaml`** — YAML frontmatter via `gopkg.in/yaml.v3`
- **`meta/toml`** — TOML frontmatter via `github.com/BurntSushi/toml`

**Template functions** (`tmpl/`):
- **`tmpl/funcs`** — stdlib-only `template.FuncMap`; `funcs.FuncMap()` returns the full map, `funcs.Merge(maps...)` combines FuncMaps. Functions follow Go stdlib argument order (subject first). See `tmpl/funcs/TODO.md` for planned additions.

**Tools and examples:**
- **`cmd/swapfrontmatter`** — CLI to convert frontmatter between YAML/JSON/Email formats; flags: `-from`, `-to`, `-write`
- **`sample/`** — complete working example: JSON frontmatter, HTML content, `text/template` macros, page templates, pretty-print, file output to `public/`

### Adding a renderer

Implement `func(io.Writer, io.Reader, data any) error` and include it in the pipeline
slice passed to `Render`. The `data` argument is the page's `ContentSourceConfig`. Use
`ssg.Must(r)` to wrap a renderer constructor that returns `(Renderer, error)`.

### Template loading

`NewPageRender()` (template.go) recursively discovers `*.html` templates under a layout
directory, keyed by filename without extension. Templates are executed with the full
`ContentSourceConfig` as the data context — `{{.Title}}`, `{{.Content}}`, etc.

**Block override constraint:** all templates in the same directory share one template
set. If two sibling templates both `{{define "main"}}`, Go's template engine errors.
Each template that overrides a block needs its own subdirectory:

```
layout/
  baseof.html        ← defines {{block "main" .}}
  post/
    index.html       ← {{define "main"}} safe — isolated set
  tag-list/
    index.html       ← {{define "main"}} safe — isolated set
```
