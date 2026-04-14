package ssg

// Plugin is the single interface for all pipeline stages: load, expand,
// contract, enrich, and materialize. Each Plugin receives the full artifact
// set and may add, remove, or modify entries.
//
//	type Plugin func(ctx *Context, artifacts *[]Artifact) error
type Plugin func(ctx *Context, artifacts *[]Artifact) error

// Artifact is a single unit of work: the page metadata plus the pipeline
// that will materialize it into output. Separating data (Meta) from behavior
// (Pipeline) keeps ContentSourceConfig a clean data map.
type Artifact struct {
	Meta     ContentSourceConfig
	Pipeline Pipeline
}
