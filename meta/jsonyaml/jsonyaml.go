// Package jsonyaml provides a MetaLoader that accepts either YAML or JSON frontmatter.
//
// YAML frontmatter is delimited by ---\n ... \n---\n:
//
//	---
//	title: My Post
//	tags: [go, web]
//	---
//	body content here
//
// JSON frontmatter is a JSON object terminated by a closing }\n on its own line:
//
//	{
//	"title": "My Post",
//	"tags": ["go", "web"]
//	}
//	body content here
//
// YAML is parsed via yamlite (no heavy YAML library dependency).
package jsonyaml

import (
	"bytes"
	stdjson "encoding/json"

	"github.com/client9/yamlite"
)

// Loader parses YAML or JSON frontmatter and returns the metadata and body.
// Files with no recognized frontmatter prefix are returned as body-only with empty metadata.
var Loader = func(raw []byte) (map[string]any, []byte, error) {
	if bytes.HasPrefix(raw, []byte("---\n")) {
		return loadYAML(raw)
	}
	if bytes.HasPrefix(raw, []byte("{\n")) {
		return loadJSON(raw)
	}
	return map[string]any{}, raw, nil
}

func loadYAML(raw []byte) (map[string]any, []byte, error) {
	prefix := []byte("---\n")
	head, body, found := bytes.Cut(raw[len(prefix):], []byte("\n---\n"))
	if !found {
		return map[string]any{}, raw, nil
	}
	jsonBytes, err := yamlite.Convert(string(head))
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]any{}
	if err := stdjson.Unmarshal(jsonBytes, &meta); err != nil {
		return nil, nil, err
	}
	return meta, body, nil
}

func loadJSON(raw []byte) (map[string]any, []byte, error) {
	sep := []byte("\n}\n")
	idx := bytes.Index(raw, sep)
	if idx == -1 {
		return map[string]any{}, raw, nil
	}
	cut := idx + len(sep)
	head, body := raw[:cut], raw[cut:]
	meta := map[string]any{}
	if err := stdjson.Unmarshal(head, &meta); err != nil {
		return nil, nil, err
	}
	return meta, body, nil
}
