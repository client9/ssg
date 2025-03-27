package ssg

import (
	"bytes"
	"strings"
	"testing"
)

func TestSplitNoHead(t *testing.T) {

	in := strings.TrimSpace(`
content
`)

	cs := ContentSplitter{}
	cs.Register(HeadYaml)
	bin := []byte(in)
	name, head, body := cs.Split(bin)
	if name != "" {
		t.Errorf("Expected empty head name, got %q", name)
	}
	if head != nil {
		t.Errorf("Expected empty head, got %q", head)
	}
	if bytes.Compare(body, bin) != 0 {
		t.Errorf("Expected body of %q, got %q", in, body)
	}
}

func TestContentSplitYaml(t *testing.T) {

	in := strings.TrimSpace(`
---
foo: bar
---
content
`)
	bin := []byte(in)
	cs := ContentSplitter{}
	cs.Register(HeadYaml)
	name, head, body := cs.Split(bin)

	if name != "yaml" {
		t.Errorf("Expected 'yaml' got %q", name)
	}
	if bytes.Compare(head, []byte("foo: bar")) != 0 {
		t.Errorf("Expected foo:bar, got %q", head)
	}
	if bytes.Compare(body, []byte("content")) != 0 {
		t.Errorf("Expected content, got %q", body)
	}
}
