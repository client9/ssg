package jsonyaml

import (
	"reflect"
	"testing"
)

func TestLoaderYAML(t *testing.T) {
	raw := []byte("---\ntitle: Hello\ntags: [go, web]\n---\nbody here\n")
	meta, body, err := Loader(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "body here\n" {
		t.Errorf("body = %q", body)
	}
	want := map[string]any{"title": "Hello", "tags": []any{"go", "web"}}
	if !reflect.DeepEqual(meta, want) {
		t.Errorf("meta = %v, want %v", meta, want)
	}
}

func TestLoaderJSON(t *testing.T) {
	raw := []byte("{\n\"title\": \"Hello\",\n\"count\": 3\n}\nbody here\n")
	meta, body, err := Loader(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "body here\n" {
		t.Errorf("body = %q", body)
	}
	if meta["title"] != "Hello" {
		t.Errorf("title = %v", meta["title"])
	}
}

func TestLoaderNoFrontmatter(t *testing.T) {
	raw := []byte("just body\n")
	meta, body, err := Loader(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "just body\n" {
		t.Errorf("body = %q", body)
	}
	if len(meta) != 0 {
		t.Errorf("expected empty meta, got %v", meta)
	}
}
