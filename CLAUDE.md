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

### Processing model

```
load inputs into memory
  ↓
enrich / expand / contract
  add or remove artifacts, derive metadata
  ↓
materialize
  each artifact runs its own pipeline → emit outputs
```

All three phases use the same `Plugin` interface. `FileWalker` + `Rules` is the only special primitive.

### Core types

```go
// Site-wide state passed to every Plugin and Stage.
type Context struct {
    Globals   map[string]any
    OutputDir string
    Logger    *log.Logger
}

// The single interface for all pipeline phases: load, expand, contract, materialize.
type Plugin func(ctx *Context, artifacts *[]Artifact) error

// A single unit of work: metadata plus the pipeline that will produce its output.
type Artifact struct {
    Meta     ContentSourceConfig
    Pipeline Pipeline
}

// A named sequence of stages executed in order.
type Pipeline struct { /* name, stages — use NewPipeline to construct */ }

func NewPipeline(name string, stages ...Stage) Pipeline

// A named pipeline step. Use Step[I,O] to create one from a typed function.
// The pipeline carries both content (typed I→O) and metadata (ContentSourceConfig).
// A step can transform content, mutate metadata, or both.
type Stage interface {
    Name() string
    Run(ctx *Context, cfg ContentSourceConfig, in any) (any, error)
}

// Step wraps a typed function into a Stage. any is confined here.
func Step[I, O any](name string, fn func(*Context, ContentSourceConfig, I) (O, error)) Stage

// Page metadata — a plain map with typed accessors:
// .OutputFile(), .TemplateName(), .InputFile(), .SourcePath(), .Get(key), .Clone()
type ContentSourceConfig map[string]any

// Parses raw file bytes into frontmatter metadata and body.
// Returning a nil map signals skip. Return type is map[string]any so
// loader implementations have no dependency on the ssg module.
type MetaLoader func(raw []byte) (map[string]any, []byte, error)

// Maps a glob pattern to a loader and a pipeline.
// Rules are tried in order; first match wins. nil Loader = skip.
type Rule struct {
    Pattern  string
    Loader   MetaLoader
    Pipeline Pipeline
}

// Maps a relative input path to a relative output path.
// Return "" to skip. Built-ins: CleanURLs, UglyURLs, SlugNormalize.
type PathTransformer func(relPath string) string
```

### The pipeline carries content AND metadata together

Each pipeline step receives both:
- **Content** — typed input value `I`, returns typed output value `O`
- **Metadata** — `ContentSourceConfig` (a `map[string]any`), mutable in place

A step can do any combination:

| What the step does | Content signature | Metadata |
|---|---|---|
| Transform content | `[]byte → []byte` | ignored |
| Mutate metadata only (pass-through) | `any → any` (returns `in` unchanged) | mutated |
| Transform content + read metadata | `[]byte → []byte` | read only |
| Terminal sink | `[]byte → struct{}` | read only |

The `any → any` pass-through works because `in.(any)` always succeeds in `Step`, so the underlying value (e.g. `[]byte`) is preserved for the next stage's type assertion.

Examples:
- `SetOutputFile` — mutates `cfg["OutputFile"]`, returns `in` unchanged (`any → any`)
- `SetTemplateName` — mutates `cfg["TemplateName"]`, returns `in` unchanged (`any → any`)
- `markdown.New()` — converts `[]byte` markdown to `[]byte` HTML, ignores `cfg`
- `minify.New()` — reads `cfg.OutputFile()` for MIME type, transforms `[]byte → []byte`
- `NewPageRender` — reads `cfg.TemplateName()`, writes `cfg["Content"]`, transforms `[]byte → []byte`
- `WriteOutput` — reads `cfg.OutputFile()` and `ctx.OutputDir`, writes to disk (`[]byte → struct{}`)
- `FanOut(name, branches...)` — runs each branch Pipeline with the same input, aggregates errors

### The three phases

1. **Load** — `FileWalker(contentDir, rules)` returns a Plugin that walks files,
   matches each against Rules, calls the MetaLoader once per file, and creates one
   `Artifact` per matched file (carrying the Rule's Pipeline).

2. **Enrich / expand / contract** — user Plugins operate on `*[]Artifact`:
   `FilterArtifacts(fn)`, `GroupByStrings` + `NewPage` for taxonomy, sort, build nav.

3. **Materialize** — `Render` merges `ctx.Globals` into each artifact's Meta (page
   wins on collision), then runs each artifact's own `Pipeline` via `RunPipeline`.

All phases are Plugins. There is no structural boundary between them:

```go
ctx := &ssg.Context{Globals: map[string]any{"Site": site}, OutputDir: "public"}

var artifacts []ssg.Artifact
for _, p := range []ssg.Plugin{
    ssg.FileWalker("content", rules),
    mysite.RemoveDrafts,
    mysite.AddTaxonomy,
    ssg.Render,
} {
    if err := p(ctx, &artifacts); err != nil { log.Fatal(err) }
}
```

### Rule and pipeline example

```go
rules := []ssg.Rule{
    {
        Pattern: "**/*.md",
        Loader:  metayaml.Loader,
        Pipeline: ssg.NewPipeline("post",
            ssg.SetOutputFile(ssg.CleanURLs(".md", ".html")), // metadata only
            ssg.SetTemplateName("post.html"),                  // metadata only
            markdown.New(),                                    // []byte → []byte
            ssg.Must(ssg.NewPageRender("layout", fns)),        // []byte → []byte, reads+writes cfg
            ssg.WriteOutput,                                   // []byte → struct{} (terminal)
        ),
    },
    {Pattern: "**/_*"}, // nil Loader: skip draft files
}
```

### Pipeline execution helpers

```go
// Run a Pipeline and return the final typed result.
RunPipeline[T any](ctx *Context, cfg ContentSourceConfig, p Pipeline, input any) (T, error)
```

### Adding a pipeline step

Implement a typed function and wrap it with `Step`:

```go
// Content-transforming step:
var MyStep = ssg.Step("my-step", func(ctx *ssg.Context, cfg ssg.ContentSourceConfig, in []byte) ([]byte, error) {
    // transform in, optionally read cfg
    return result, nil
})

// Metadata-only step (pass-through):
func SetFoo(val string) ssg.Stage {
    return ssg.Step("set-foo", func(_ *ssg.Context, cfg ssg.ContentSourceConfig, in any) (any, error) {
        cfg["Foo"] = val
        return in, nil
    })
}
```

Use `ssg.Must(ssg.NewPageRender("layout", fns))` to inline constructors that return `(Stage, error)`.

### One-to-many outputs

Use `FanOut` inside a pipeline to produce multiple output files from one source.
Each branch is a full Pipeline; all branches receive the same input:

```go
Pipeline: ssg.NewPipeline("post",
    ssg.FanOut("outputs",
        ssg.NewPipeline("html", ssg.SetOutputFile(ssg.CleanURLs(".md", ".html")), markdown.New(), ssg.WriteOutput),
        ssg.NewPipeline("txt",  ssg.SetOutputFile(ssg.UglyURLs(".md", ".txt")),  plaintext.New(), ssg.WriteOutput),
    ),
),
```

### Built-in MetaLoaders

- `ssg.Passthrough` — returns raw bytes as body with empty metadata; use for assets
- `ssg.Skip` — unconditionally skips the file (same as nil Loader, but explicit)

### Synthetic pages

`NewPage(outputFile, templateName, data, pipeline)` creates an Artifact not backed
by a file. Used for taxonomy indexes, RSS feeds, etc.:

```go
*artifacts = append(*artifacts, ssg.NewPage(
    "tags/go/index.html", "tag-list/index.html",
    map[string]any{"Tag": "go", "Pages": tagMetas},
    ssg.NewPipeline("tag",
        ssg.Must(ssg.NewPageRender("layout", fns)),
        ssg.WriteOutput,
    ),
))
```

### Taxonomy helpers

`taxonomy.go`:
- `GroupByString(artifacts, field)` — group by a single string field (e.g. `"Category"`)
- `GroupByStrings(artifacts, field)` — group by a multi-value field (e.g. `"Tags"`);
  handles `[]string`, `[]any`, and bare `string`

### Filtering

`filter.go`:
- `FilterArtifacts(fn) Plugin` — returns a Plugin that retains artifacts where `fn` returns true

### ContentSourceConfig

`map[string]any` with typed accessors. The map is the right abstraction: system fields
are type-safe via accessors, user frontmatter is dynamic, templates access all keys
uniformly via `{{.Title}}` etc. Storing the `Pipeline` on a separate `Artifact`
struct keeps the map clean.

Known system keys: `OutputFile`, `TemplateName`, `InputFile`, `SourcePath`, `Content`.

### Sub-packages

Each is a separate Go module in `go.work`. Meta sub-modules have **no dependency on
`github.com/client9/ssg`**. Render sub-modules implement `ssg.DynStage` and do depend on ssg.

**Pipeline stages** (`render/`) — return `ssg.Stage`:
- **`render/htmlclean`** — normalizes HTML fragments via `golang.org/x/net/html`
- **`render/markdown`** — Goldmark Markdown→HTML; `markdown.New()` or `markdown.NewGoldmark(g)`
- **`render/minify`** — minifies HTML/CSS/JS/SVG; MIME type from output file extension
- **`render/shortcode`** — `$cmd[args]{body}` macro expansion engine

**Metadata loaders** (`meta/`) — each exports `var Loader MetaLoader`:
- **`meta/yaml`** — YAML frontmatter (`---\n...\n---\n`)
- **`meta/toml`** — TOML frontmatter (`+++\n...\n+++\n`)
- **`meta/json`** — JSON frontmatter (`{\n...\n}\n`)
- **`meta/email`** — email-style `Key: Value` headers; `email.NewLoader(transformers...)`

**Template functions** (`tmpl/`):
- **`tmpl/stdfuncs`** — stdlib-only `template.FuncMap`; `stdfuncs.FuncMap()`, `stdfuncs.Merge(...)`

### Template loading

`NewPageRender(tdir, fns)` discovers `*.html` templates under a layout directory.
Template selection uses `cfg.TemplateName()` — directory portion routes to the right
set, filename selects the template within it.

**Block override constraint:** templates in the same directory share one set. If two
siblings both `{{define "main"}}`, Go's template engine errors. Each block-overriding
template needs its own subdirectory:

```
layout/
  baseof.html        ← defines {{block "main" .}}
  post/
    index.html       ← {{define "main"}} safe — isolated set
  tag-list/
    index.html       ← {{define "main"}} safe — isolated set
```
