// Package yaml provides a MetaLoader for YAML-formatted frontmatter.
//
// Frontmatter is delimited by ---\n ... \n---\n:
//
//	---
//	title: My Post
//	tags: [go, web]
//	---
//	body content here
package yaml

import (
	"bytes"

	goyaml "go.yaml.in/yaml/v4"
)

// Loader parses YAML frontmatter and returns the metadata and body.
// Files with no ---\n prefix are returned as body-only with empty metadata.
var Loader = func(raw []byte) (map[string]any, []byte, error) {
	head, body := split(raw)
	if head == nil {
		return map[string]any{}, body, nil
	}
	meta := map[string]any{}
	if err := goyaml.Unmarshal(head, &meta); err != nil {
		return nil, nil, err
	}
	return meta, body, nil
}

func split(raw []byte) (head, body []byte) {
	prefix := []byte("---\n")
	if !bytes.HasPrefix(raw, prefix) {
		return nil, raw
	}
	head, body, found := bytes.Cut(raw[len(prefix):], []byte("\n---\n"))
	if !found {
		return nil, raw
	}
	return head, body
}
