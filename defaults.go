package ssg

import "maps"

// LoadDefaults fills in zero-value fields of conf with sensible defaults.
func LoadDefaults(conf *LoadConfig) {
	if conf.ContentDir == "" {
		conf.ContentDir = "content"
	}
}

// NewPage creates a ContentSourceConfig with the given output path, template,
// and data. It is the standard constructor for all pages — both those loaded
// from files (via LoadContent) and those created programmatically.
//
// The data map is merged first; OutputFile and TemplateName are then set from
// the explicit arguments and always win. Content is initialised to an empty
// byte slice; LoadContent overwrites it with the file body.
//
//	p := ssg.NewPage("tags/go/index.html", "tag-list/index.html", map[string]any{
//	    "Tag":   "go",
//	    "Pages": tagPages,
//	})
func NewPage(outputFile, templateName string, data map[string]any) ContentSourceConfig {
	p := make(ContentSourceConfig, len(data)+3)
	maps.Copy(p, data)
	p["OutputFile"] = outputFile
	p["TemplateName"] = templateName
	p["Content"] = []byte{}
	return p
}

// ContentSourceConfig holds metadata and content for a single page.
// It is a plain map so that frontmatter parsers, synthetic page builders,
// and templates can all read and write it without a fixed schema.
type ContentSourceConfig map[string]any

func (csc ContentSourceConfig) TemplateName() string {
	if val, ok := csc["TemplateName"]; ok {
		return val.(string)
	}
	return ""
}

func (csc ContentSourceConfig) OutputFile() string {
	if val, ok := csc["OutputFile"]; ok {
		return val.(string)
	}
	return ""
}

func (csc ContentSourceConfig) InputFile() string {
	if val, ok := csc["InputFile"]; ok {
		return val.(string)
	}
	return ""
}

func (csc ContentSourceConfig) Get(key string) string {
	if val, ok := csc[key]; ok {
		if sval, ok := val.(string); ok {
			return sval
		}
	}
	return ""
}
