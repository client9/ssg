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
	bin := []byte(in)
	head, body := Splitter(HeadYaml, bin)
	if head != nil {
		t.Errorf("Expected empty head, got %q", head)
	}
	if bytes.Compare(body, bin) != 0 {
		t.Errorf("Expected body of %q, got %q", in, body)
	}
}
func TestContentSplitJson(t *testing.T) {
	in := strings.TrimSpace(`
{
"foo": "bar"
}
content
`)
	bin := []byte(in)
	head, body := Splitter(HeadJson, bin)

	want := []byte("{\n\"foo\": \"bar\"\n}\n")
	if bytes.Compare(want, head) != 0 {
		t.Errorf("Expected %q,  got %q", want, head)
	}
	if bytes.Compare(body, []byte("content")) != 0 {
		t.Errorf("Expected content, got %q", body)
	}

	// make sure it parses
	_, err := ParseMetaJson(head)
	if err != nil {
		t.Errorf("Json didn't parse,%v", err)
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
	head, body := Splitter(HeadYaml, bin)

	if bytes.Compare(head, []byte("foo: bar")) != 0 {
		t.Errorf("Expected foo:bar, got %q", head)
	}
	if bytes.Compare(body, []byte("content")) != 0 {
		t.Errorf("Expected content, got %q", body)
	}
}
