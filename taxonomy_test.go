package ssg

import "testing"

func artifactWithMeta(meta map[string]any) Artifact {
	return Artifact{Meta: ContentSourceConfig(meta)}
}

func TestGroupByString(t *testing.T) {
	artifacts := []Artifact{
		artifactWithMeta(map[string]any{"Category": "go", "Title": "A"}),
		artifactWithMeta(map[string]any{"Category": "go", "Title": "B"}),
		artifactWithMeta(map[string]any{"Category": "web", "Title": "C"}),
		artifactWithMeta(map[string]any{"Title": "D"}),                 // no Category
		artifactWithMeta(map[string]any{"Category": "", "Title": "E"}), // empty Category
	}

	got := GroupByString(artifacts, "Category")

	if len(got["go"]) != 2 {
		t.Errorf("expected 2 in 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 1 {
		t.Errorf("expected 1 in 'web', got %d", len(got["web"]))
	}
	if _, ok := got[""]; ok {
		t.Error("empty string should not be a key")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 groups, got %d", len(got))
	}
}

func TestGroupByStrings_sliceString(t *testing.T) {
	artifacts := []Artifact{
		artifactWithMeta(map[string]any{"Tags": []string{"go", "web"}, "Title": "A"}),
		artifactWithMeta(map[string]any{"Tags": []string{"go"}, "Title": "B"}),
		artifactWithMeta(map[string]any{"Tags": []string{"web"}, "Title": "C"}),
	}

	got := GroupByStrings(artifacts, "Tags")

	if len(got["go"]) != 2 {
		t.Errorf("expected 2 tagged 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 2 {
		t.Errorf("expected 2 tagged 'web', got %d", len(got["web"]))
	}
}

func TestGroupByStrings_sliceAny(t *testing.T) {
	artifacts := []Artifact{
		artifactWithMeta(map[string]any{"Tags": []any{"go", "web"}, "Title": "A"}),
		artifactWithMeta(map[string]any{"Tags": []any{"go"}, "Title": "B"}),
	}

	got := GroupByStrings(artifacts, "Tags")

	if len(got["go"]) != 2 {
		t.Errorf("expected 2 tagged 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 1 {
		t.Errorf("expected 1 tagged 'web', got %d", len(got["web"]))
	}
}

func TestGroupByStrings_bareString(t *testing.T) {
	artifacts := []Artifact{
		artifactWithMeta(map[string]any{"Tags": "go", "Title": "A"}),
		artifactWithMeta(map[string]any{"Tags": "web", "Title": "B"}),
	}

	got := GroupByStrings(artifacts, "Tags")

	if len(got["go"]) != 1 {
		t.Errorf("expected 1 tagged 'go', got %d", len(got["go"]))
	}
	if len(got["web"]) != 1 {
		t.Errorf("expected 1 tagged 'web', got %d", len(got["web"]))
	}
}

func TestGroupByStrings_missingField(t *testing.T) {
	artifacts := []Artifact{artifactWithMeta(map[string]any{"Title": "A"})}
	got := GroupByStrings(artifacts, "Tags")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}

func TestNewPage(t *testing.T) {
	pipelineNop := Pipeline{name: "nop", stages: nil}
	data := map[string]any{"Tag": "go"}
	a := NewPage("tags/go/index.html", "tag-list.html", data, pipelineNop)

	if got := a.Meta.OutputFile(); got != "tags/go/index.html" {
		t.Errorf("OutputFile = %q", got)
	}
	if got := a.Meta.TemplateName(); got != "tag-list.html" {
		t.Errorf("TemplateName = %q", got)
	}
	if _, ok := a.Meta["Content"].([]byte); !ok {
		t.Error("Content should be []byte")
	}
	if got := a.Meta["Tag"]; got != "go" {
		t.Errorf("Tag = %v", got)
	}
}

func TestNewPage_dataCannotOverrideOutputFile(t *testing.T) {
	pipelineNop := Pipeline{name: "nop", stages: nil}
	data := map[string]any{"OutputFile": "should-be-ignored.html"}
	a := NewPage("real/index.html", "base.html", data, pipelineNop)
	if got := a.Meta.OutputFile(); got != "real/index.html" {
		t.Errorf("OutputFile = %q, want real/index.html", got)
	}
}

func TestNewPage_nilData(t *testing.T) {
	pipelineNop := Pipeline{name: "nop", stages: nil}
	a := NewPage("out.html", "base.html", nil, pipelineNop)
	if a.Meta.OutputFile() != "out.html" {
		t.Error("OutputFile not set with nil data")
	}
}
