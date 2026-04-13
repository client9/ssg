// Package json provides a MetaLoader for JSON-formatted frontmatter.
//
// Frontmatter is a JSON object at the start of the file, terminated by
// a closing }\n on its own line:
//
//	{
//	"title": "My Post",
//	"tags": ["go", "web"]
//	}
//	body content here
package json

import (
	"bytes"
	stdjson "encoding/json"
)

// Loader parses JSON frontmatter and returns the metadata and body.
// Files with no {\n prefix are returned as body-only with empty metadata.
var Loader = func(raw []byte) (map[string]any, []byte, error) {
	head, body := split(raw)
	if head == nil {
		return map[string]any{}, body, nil
	}
	meta := map[string]any{}
	if err := stdjson.Unmarshal(head, &meta); err != nil {
		return nil, nil, err
	}
	return meta, body, nil
}

// split separates a JSON frontmatter block from the body. The head includes
// the opening { and closing } so that json.Unmarshal can parse it directly.
func split(raw []byte) (head, body []byte) {
	prefix := []byte("{\n")
	sep := []byte("\n}\n")
	if !bytes.HasPrefix(raw, prefix) {
		return nil, raw
	}
	idx := bytes.Index(raw, sep)
	if idx == -1 {
		return nil, raw
	}
	cut := idx + len(sep)
	return raw[:cut], raw[cut:]
}
