# shortcode — Claude guidance

## Project
Go module `github.com/client9/ssg/render/shortcode`. A plain-text embedded macro/template engine. Go 1.23+.

## File layout
One concept per file. Test file alongside implementation (`foo.go` / `foo_test.go`).
Current files:
- `macro.go` — core types, parser, Render, RenderDocument, PositionError
- `named.go` — NamedFunc, MakeNamed, RegisterNamed
- `entity.go` — Entity handler
- `csvtable.go` — CSVTable handler and macro-aware CSV parser

## Conventions

**HandlerFunc signature is stable — do not change it.**
```go
type HandlerFunc func(ctx *Context, args []string, body string) string
```

**Error handling is accumulation, not propagation.**
Handlers call `ctx.AddError(err)` and return a string (empty or a placeholder).
Never change Render to return an error. RenderDocument is the top-level entry point.

**Named args live in []string as "key=value" strings.**
ParseNamedArgs and MakeNamed convert them. Do not add a map parameter to HandlerFunc.

**Recursive rendering is explicit and opt-in.**
Handlers call `ctx.Render(body)` themselves. The engine never auto-renders bodies.

**No dependencies beyond tablewriter for built-in handlers.**
Core parsing and rendering has no external dependencies.

## Testing
Table-driven tests where there are multiple input/output cases.
Test files are in the same package (not `shortcode_test`).
Run: `go test ./...`
