package ssg

import "testing"

func TestGroupByString(t *testing.T) {
	pages := []ContentSourceConfig{
		{"Category": "go", "Title": "A"},
		{"Category": "go", "Title": "B"},
		{"Category": "web", "Title": "C"},
		{"Title": "D"},           // no Category
		{"Category": "", "Title": "E"}, // empty Category
	}

	got := GroupByString(pages, "Category")

	if len(got["go"]) != 2 {
		t.Errorf("expected 2 pages in 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 1 {
		t.Errorf("expected 1 page in 'web', got %d", len(got["web"]))
	}
	if _, ok := got[""]; ok {
		t.Error("empty string should not be a key")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 groups, got %d", len(got))
	}
}

func TestGroupByStrings_sliceString(t *testing.T) {
	pages := []ContentSourceConfig{
		{"Tags": []string{"go", "web"}, "Title": "A"},
		{"Tags": []string{"go"}, "Title": "B"},
		{"Tags": []string{"web"}, "Title": "C"},
	}

	got := GroupByStrings(pages, "Tags")

	if len(got["go"]) != 2 {
		t.Errorf("expected 2 pages tagged 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 2 {
		t.Errorf("expected 2 pages tagged 'web', got %d", len(got["web"]))
	}
}

func TestGroupByStrings_sliceAny(t *testing.T) {
	// []any is what YAML/JSON parsers typically produce
	pages := []ContentSourceConfig{
		{"Tags": []any{"go", "web"}, "Title": "A"},
		{"Tags": []any{"go"}, "Title": "B"},
	}

	got := GroupByStrings(pages, "Tags")

	if len(got["go"]) != 2 {
		t.Errorf("expected 2 pages tagged 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 1 {
		t.Errorf("expected 1 page tagged 'web', got %d", len(got["web"]))
	}
}

func TestGroupByStrings_bareString(t *testing.T) {
	// single tag written as a plain string, not a list
	pages := []ContentSourceConfig{
		{"Tags": "go", "Title": "A"},
		{"Tags": "web", "Title": "B"},
	}

	got := GroupByStrings(pages, "Tags")

	if len(got["go"]) != 1 {
		t.Errorf("expected 1 page tagged 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 1 {
		t.Errorf("expected 1 page tagged 'web', got %d", len(got["web"]))
	}
}

func TestGroupByStrings_missingField(t *testing.T) {
	pages := []ContentSourceConfig{
		{"Title": "A"}, // no Tags field
	}
	got := GroupByStrings(pages, "Tags")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}

func TestNewPage(t *testing.T) {
	data := map[string]any{"Tag": "go", "Pages": []ContentSourceConfig{}}
	p := NewPage("tags/go/index.html", "tag-list.html", data)

	if got := p.OutputFile(); got != "tags/go/index.html" {
		t.Errorf("OutputFile = %q", got)
	}
	if got := p.TemplateName(); got != "tag-list.html" {
		t.Errorf("TemplateName = %q", got)
	}
	if _, ok := p["Content"].([]byte); !ok {
		t.Error("Content should be []byte")
	}
	if got := p["Tag"]; got != "go" {
		t.Errorf("Tag = %v", got)
	}
}

func TestNewPage_dataCannotOverrideOutputFile(t *testing.T) {
	// OutputFile in data must not win over the explicit argument
	data := map[string]any{"OutputFile": "should-be-ignored.html"}
	p := NewPage("real/index.html", "base.html", data)

	if got := p.OutputFile(); got != "real/index.html" {
		t.Errorf("OutputFile = %q, want real/index.html", got)
	}
}

func TestNewPage_nilData(t *testing.T) {
	p := NewPage("out.html", "base.html", nil)
	if p.OutputFile() != "out.html" {
		t.Error("OutputFile not set with nil data")
	}
}
