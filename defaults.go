package ssg

import "maps"

// NewPage creates an Artifact for a page not backed by a file — taxonomy indexes,
// tag listings, RSS feeds, or any other programmatically-generated output.
//
// data is merged first; OutputFile and TemplateName are then set from the
// explicit arguments and always win. Content is initialised to an empty byte
// slice so the pipeline has something to read.
//
//	p := ssg.NewPage(
//	    "tags/go/index.html", "tag-list/index.html",
//	    map[string]any{"Tag": "go", "Pages": tagPages},
//	    pipeline,
//	)
func NewPage(outputFile, templateName string, data map[string]any, pipeline Pipeline) Artifact {
	meta := make(ContentSourceConfig, len(data)+3)
	maps.Copy(meta, data)
	meta["OutputFile"] = outputFile
	meta["TemplateName"] = templateName
	meta["Content"] = []byte{}
	return Artifact{Meta: meta, Pipeline: pipeline}
}

// ContentSourceConfig holds metadata and content for a single page.
// It is a plain map so that frontmatter parsers, synthetic page builders,
// and templates can all read and write it without a fixed schema.
type ContentSourceConfig map[string]any

// Clone returns a shallow copy of the map.
func (csc ContentSourceConfig) Clone() ContentSourceConfig {
	return ContentSourceConfig(maps.Clone(map[string]any(csc)))
}

func (csc ContentSourceConfig) TemplateName() string {
	v, _ := csc["TemplateName"].(string)
	return v
}

func (csc ContentSourceConfig) OutputFile() string {
	v, _ := csc["OutputFile"].(string)
	return v
}

func (csc ContentSourceConfig) InputFile() string {
	v, _ := csc["InputFile"].(string)
	return v
}

// SourcePath returns the path of the source file relative to ContentDir.
// Used by SetOutputFile to derive the output path via a PathTransformer.
func (csc ContentSourceConfig) SourcePath() string {
	v, _ := csc["SourcePath"].(string)
	return v
}

func (csc ContentSourceConfig) Get(key string) string {
	v, _ := csc[key].(string)
	return v
}
