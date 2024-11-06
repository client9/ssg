package ssg

import (
	"strings"
	"testing"
)

func TestSplitNoHead(t *testing.T) {

	in := strings.TrimSpace(`
content
`)

	cs := ContentSplitter{}
	cs.Register(HeadYaml)
	name, head, body := cs.Split(in)
	if name != "" {
		t.Errorf("Expected empty head name, got %q", name)
	}
	if head != "" {
		t.Errorf("Expected empty head, got %q", head)
	}
	if body != in {
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
	cs := ContentSplitter{}
	cs.Register(HeadYaml)
	name, head, body := cs.Split(in)

	if name != "yaml" {
		t.Errorf("Expected 'yaml' got %q", name)
	}
	if head != "foo: bar" {
		t.Errorf("Expected foo:bar, got %q", head)
	}
	if body != "content" {
		t.Errorf("Expected content, got %q", body)
	}
}
