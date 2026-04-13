// Package toml provides a MetaLoader for TOML-formatted frontmatter.
//
// Frontmatter is delimited by +++\n ... \n+++\n:
//
//	+++
//	title = "My Post"
//	tags = ["go", "web"]
//	+++
//	body content here
package toml

import (
	"bytes"

	btoml "github.com/BurntSushi/toml"
	"github.com/client9/ssg"
)

// Loader parses TOML frontmatter and returns the metadata and body.
// Files with no +++\n prefix are returned as body-only with empty metadata.
var Loader ssg.MetaLoader = func(raw []byte) (ssg.ContentSourceConfig, []byte, error) {
	head, body := split(raw)
	if head == nil {
		return ssg.ContentSourceConfig{}, body, nil
	}
	meta := ssg.ContentSourceConfig{}
	if _, err := btoml.NewDecoder(bytes.NewReader(head)).Decode(&meta); err != nil {
		return nil, nil, err
	}
	return meta, body, nil
}

func split(raw []byte) (head, body []byte) {
	prefix := []byte("+++\n")
	if !bytes.HasPrefix(raw, prefix) {
		return nil, raw
	}
	head, body, found := bytes.Cut(raw[len(prefix):], []byte("\n+++\n"))
	if !found {
		return nil, raw
	}
	return head, body
}
