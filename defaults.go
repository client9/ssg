package ssg

func ConfigDefault(config *SiteConfig) {
	if config.ContentDir == "" {
		config.ContentDir = "content"
	}
	if config.MetaSplit == nil {
		config.MetaSplit = MetaSplitJson
	}
	if config.MetaParser == nil {
		config.MetaParser = MetaParseJson
	}
	if config.OutputExt == "" {
		config.OutputExt = ".html"
	}
	if config.InputExt == "" {
		config.InputExt = ".html"
	}
	if config.IndexSource == "" {
		config.IndexSource = "index.html"
	}
	if config.IndexDest == "" {
		config.IndexDest = "index.html"
	}
	if config.BaseTemplate == "" {
		config.BaseTemplate = "baseof.html"
	}
}

// Default ContentSource -- it's just a map any with some specific
// accessors
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
