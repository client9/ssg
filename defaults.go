package ssg

// LoadDefaults fills in zero-value fields of conf with sensible defaults
// for a standard HTML site.
func LoadDefaults(conf *LoadConfig) {
	if conf.ContentDir == "" {
		conf.ContentDir = "content"
	}
	if conf.MetaSplit == nil {
		conf.MetaSplit = MetaSplitJson
	}
	if conf.MetaParser == nil {
		conf.MetaParser = MetaParseJson
	}
	if conf.OutputExt == "" {
		conf.OutputExt = ".html"
	}
	if conf.InputExt == "" {
		conf.InputExt = ".html"
	}
	if conf.IndexSource == "" {
		conf.IndexSource = "index.html"
	}
	if conf.IndexDest == "" {
		conf.IndexDest = "index.html"
	}
	if conf.BaseTemplate == "" {
		conf.BaseTemplate = "baseof.html"
	}
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
