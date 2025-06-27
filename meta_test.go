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
	head, body := Splitter(MetaHeadYaml, bin)
	if head != nil {
		t.Errorf("Expected empty head, got %q", head)
	}
	if !bytes.Equal(body, bin) {
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
	head, body := Splitter(MetaHeadJson, bin)

	want := []byte("{\n\"foo\": \"bar\"\n}\n")
	if !bytes.Equal(want, head) {
		t.Errorf("Expected %q,  got %q", want, head)
	}
	if !bytes.Equal(body, []byte("content")) {
		t.Errorf("Expected content, got %q", body)
	}

	// make sure it parses
	_, err := MetaParseJson(head)
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
	head, body := Splitter(MetaHeadYaml, bin)

	if !bytes.Equal(head, []byte("foo: bar")) {
		t.Errorf("Expected foo:bar, got %q", head)
	}
	if !bytes.Equal(body, []byte("content")) {
		t.Errorf("Expected content, got %q", body)
	}
}
