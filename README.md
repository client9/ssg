# ssg

A document transformation pipeline and composable static site generator toolkit for Go.

`ssg` provides the building blocks for a content pipeline in three phases:

```
load inputs into memory
  ↓
enrich / expand / contract
  add or remove artifacts, derive metadata
  ↓
materialize
  each artifact runs its own pipeline → emit outputs
```

## Why? And Alternatives

There are many static site generators!

Most focus on a "no programming required" model which limits extensibility.

The closest model is [metalsmith](https://metalsmith.io).

## Core concepts

**`Context`** — site-wide state available to every Plugin and pipeline stage:

```go
type Context struct {
    Globals   map[string]any
    OutputDir string
    Logger    *log.Logger
}
```

**`Plugin`** — the single interface for all pipeline phases:

```go
type Plugin func(ctx *Context, artifacts *[]Artifact) error
```

Load, filter, expand, and materialize are all Plugins operating on the same
artifact set. No structural boundary between phases.

**`Artifact`** — one unit of work: metadata plus the pipeline that produces it:

```go
type Artifact struct {
    Meta     ContentSourceConfig // map[string]any with typed accessors
    Pipeline Pipeline
}
```

**`Pipeline`** — a named sequence of stages. Construct with `NewPipeline`:

```go
func NewPipeline(name string, stages ...Stage) Pipeline
```

**`Stage`** — a single named pipeline step. Each step receives both the current
content value and the page metadata, and can transform either or both:

```go
type Stage interface {
    Name() string
    Run(ctx *Context, cfg ContentSourceConfig, in any) (any, error)
}
```

Use `Step[I, O]` to create a `Stage` from a typed function:

```go
func Step[I, O any](name string, fn func(*Context, ContentSourceConfig, I) (O, error)) Stage
```

The pipeline carries **content and metadata together**. A step can:
- Transform content only — `[]byte → []byte`, ignore `cfg`
- Mutate metadata only — `any → any` pass-through, write to `cfg`
- Read metadata to transform content — e.g. pick a MIME type from `cfg.OutputFile()`
- Read and write metadata while transforming content — e.g. wrap rendered HTML in a layout template

**`MetaLoader`** — parses raw file bytes into frontmatter metadata and body:

```go
type MetaLoader func(raw []byte) (map[string]any, []byte, error)
```

Returning a nil map signals skip. The return type is `map[string]any` so loader
implementations have no dependency on this module.

**`Rule`** — pairs a [doublestar](https://github.com/bmatcuk/doublestar) glob pattern
with a loader and a pipeline:

```go
type Rule struct {
    Pattern  string
    Loader   MetaLoader // nil or ssg.Skip = skip without reading
    Pipeline Pipeline
}
```

## Usage

```go
ctx := &ssg.Context{
    Globals:   map[string]any{"Site": siteConfig},
    OutputDir: "public",
    Logger:    log.Default(),
}

rules := []ssg.Rule{
    {
        Pattern: "**/*.md",
        Loader:  metayaml.Loader,
        Pipeline: ssg.NewPipeline("post",
            ssg.SetOutputFile(ssg.CleanURLs(".md", ".html")), // metadata only
            ssg.SetTemplateName("post.html"),                  // metadata only
            markdown.New(),                                    // []byte → []byte
            ssg.Must(ssg.NewPageRender("layout", fns)),        // reads+writes cfg, []byte → []byte
            ssg.WriteOutput,                                   // reads cfg, terminal sink
        ),
    },
    {Pattern: "**/_*"}, // nil Loader: skip draft files
}

var artifacts []ssg.Artifact
for _, p := range []ssg.Plugin{
    ssg.FileWalker("content", rules), // Phase 1: load
    removeDrafts,                     // Phase 2: contract
    addTaxonomy,                      // Phase 2: expand
    ssg.Render,                       // Phase 3: materialize
} {
    if err := p(ctx, &artifacts); err != nil {
        log.Fatal(err)
    }
}
```

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

### Writing a pipeline step

Implement a typed function and wrap it with `Step`. The function receives both the
current content value and the mutable metadata map:

```go
// Content-transforming step (metadata ignored):
var UpperCase = ssg.Step("uppercase", func(_ *ssg.Context, _ ssg.ContentSourceConfig, in []byte) ([]byte, error) {
    return bytes.ToUpper(in), nil
})

// Metadata-only step (content passed through unchanged):
func SetCanonical(base string) ssg.Stage {
    return ssg.Step("set-canonical", func(_ *ssg.Context, cfg ssg.ContentSourceConfig, in any) (any, error) {
        cfg["Canonical"] = base + cfg.OutputFile()
        return in, nil
    })
}

// Step that reads metadata to transform content:
var AddTitle = ssg.Step("add-title", func(_ *ssg.Context, cfg ssg.ContentSourceConfig, in []byte) ([]byte, error) {
    title := cfg.Get("Title")
    return append([]byte("<h1>"+title+"</h1>\n"), in...), nil
})
```

Use `ssg.Must(ssg.NewPageRender("layout", fns))` to inline constructors that return
`(Stage, error)`.

### Filtering

```go
ssg.FilterArtifacts(func(meta ssg.ContentSourceConfig) bool {
    draft, _ := meta["draft"].(bool)
    return !draft
})
```

### Taxonomy pages

```go
byTag := ssg.GroupByStrings(artifacts, "Tags")
for tag, tagArtifacts := range byTag {
    artifacts = append(artifacts, ssg.NewPage(
        "tags/"+slug(tag)+"/index.html", "tag-list/index.html",
        map[string]any{"Tag": tag, "Pages": metaSlice(tagArtifacts)},
        tagPipeline,
    ))
}
```

### Built-in metadata steps

| Step | What it does |
|---|---|
| `SetOutputFile(transform)` | Applies a `PathTransformer` to `SourcePath`, writes `OutputFile` to cfg |
| `SetTemplateName(name)` | Writes `TemplateName` to cfg if not already set by frontmatter |

Both pass content through unchanged (`any → any`).

### Path transformers

| Function | Example |
|---|---|
| `CleanURLs(".md", ".html")` | `posts/foo.md` → `posts/foo/index.html` |
| `UglyURLs(".md", ".html")` | `posts/foo.md` → `posts/foo.html` |
| `SlugNormalize(next)` | lowercases and hyphenates before applying next |

## Sub-modules

Each sub-module is a separate Go module and can be imported independently.
Meta sub-modules have no dependency on `github.com/client9/ssg`.

### Pipeline stages (`render/`)

Each package returns a `ssg.Stage` (or a constructor for one).

| Module | Import path | Description |
|---|---|---|
| **htmlclean** | `github.com/client9/ssg/render/htmlclean` | Normalizes HTML fragments via `golang.org/x/net/html` |
| **markdown** | `github.com/client9/ssg/render/markdown` | Markdown → HTML via Goldmark with GFM and auto heading IDs |
| **minify** | `github.com/client9/ssg/render/minify` | Minifies HTML/CSS/JS/SVG; MIME type from `cfg.OutputFile()` |
| **shortcode** | `github.com/client9/ssg/render/shortcode` | Embedded `$cmd[args]{body}` macro engine |

The shortcode syntax:

```
$cmd
$cmd[arg1 arg2]
$cmd[name=value key="val"]
$cmd{body}
$cmd[args]{body}
$$   →  literal $
```

### Metadata loaders (`meta/`)

Each package exports a single `var Loader` of type `func([]byte) (map[string]any, []byte, error)`.

| Module | Import path | Description |
|---|---|---|
| **json** | `github.com/client9/ssg/meta/json` | JSON object frontmatter (`{\n...\n}\n`) |
| **yaml** | `github.com/client9/ssg/meta/yaml` | YAML frontmatter (`---\n...\n---\n`) via `go.yaml.in/yaml/v4` |
| **toml** | `github.com/client9/ssg/meta/toml` | TOML frontmatter (`+++\n...\n+++\n`) via `github.com/BurntSushi/toml` |
| **email** | `github.com/client9/ssg/meta/email` | Email-style `Key: Value` headers; `email.NewLoader(transformers...)` for type coercion |

The root module also provides two built-in loaders:
- `ssg.Passthrough` — returns raw bytes as body with empty metadata; use for assets
- `ssg.Skip` — unconditionally skips the file; explicit alternative to a nil `Rule.Loader`

### Template functions (`tmpl/`)

| Module | Import path | Description |
|---|---|---|
| **stdfuncs** | `github.com/client9/ssg/tmpl/stdfuncs` | Stdlib-only `template.FuncMap`; covers strings, math, collections, path, time, encoding, and URL helpers |

```go
t := template.New("page").Funcs(stdfuncs.FuncMap())

// Combine with your own:
fns := stdfuncs.Merge(stdfuncs.FuncMap(), template.FuncMap{"myFunc": myFunc})
```

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
