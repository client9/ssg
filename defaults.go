package ssg

import (
	"encoding/json"
)

func ConfigDefault(config *SiteConfig) {
	if config.OutputDir == "" {
		config.OutputDir = "public"
	}
	if config.ContentDir == "" {
		config.ContentDir = "content"
	}
	if config.TemplateDir == "" {
		config.TemplateDir = "layout"
	}
	if config.Split == nil {
		config.Split = ContentSplitJson
	}
	if config.Metaparser == nil {
		config.Metaparser = ParseMetaJson
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

// Splits a file in front/body based on JSON front matter signature.
func ContentSplitJson(s []byte) ([]byte, []byte) {
	return Splitter(HeadJson, s)
}

// ParseMetaJson is a default parser, that reads front matter as
// JSON and returns a map[string]any type.
func ParseMetaJson(s []byte) (ContentSourceConfig, error) {
	meta := ContentSourceConfig{}
	if err := json.Unmarshal(s, &meta); err != nil {
		return nil, err
	}
	return meta, nil
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
