package ssg

import (
	"testing"
)

func makeArtifact(meta map[string]any) Artifact {
	return Artifact{Meta: ContentSourceConfig(meta)}
}

func runFilter(artifacts []Artifact, fn func(ContentSourceConfig) bool) []Artifact {
	plugin := FilterArtifacts(fn)
	_ = plugin(nil, &artifacts)
	return artifacts
}

func TestFilterArtifacts_keepAll(t *testing.T) {
	artifacts := []Artifact{makeArtifact(map[string]any{"Title": "A"}), makeArtifact(map[string]any{"Title": "B"})}
	got := runFilter(artifacts, func(ContentSourceConfig) bool { return true })
	if len(got) != 2 {
		t.Errorf("expected 2 artifacts, got %d", len(got))
	}
}

func TestFilterArtifacts_removeAll(t *testing.T) {
	artifacts := []Artifact{makeArtifact(map[string]any{"Title": "A"}), makeArtifact(map[string]any{"Title": "B"})}
	got := runFilter(artifacts, func(ContentSourceConfig) bool { return false })
	if len(got) != 0 {
		t.Errorf("expected 0 artifacts, got %d", len(got))
	}
}

func TestFilterArtifacts_draftExclusion(t *testing.T) {
	artifacts := []Artifact{
		makeArtifact(map[string]any{"Title": "Published"}),
		makeArtifact(map[string]any{"Title": "Draft", "draft": true}),
		makeArtifact(map[string]any{"Title": "Also published", "draft": false}),
	}
	got := runFilter(artifacts, func(meta ContentSourceConfig) bool {
		draft, _ := meta["draft"].(bool)
		return !draft
	})
	if len(got) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(got))
	}
	for _, a := range got {
		if a.Meta["Title"] == "Draft" {
			t.Error("draft artifact should have been excluded")
		}
	}
}

func TestFilterArtifacts_byStringField(t *testing.T) {
	artifacts := []Artifact{
		makeArtifact(map[string]any{"section": "blog"}),
		makeArtifact(map[string]any{"section": "docs"}),
		makeArtifact(map[string]any{"section": "blog"}),
	}
	got := runFilter(artifacts, func(meta ContentSourceConfig) bool {
		return meta.Get("section") == "blog"
	})
	if len(got) != 2 {
		t.Errorf("expected 2 blog artifacts, got %d", len(got))
	}
}

func TestFilterArtifacts_empty(t *testing.T) {
	got := runFilter(nil, func(ContentSourceConfig) bool { return true })
	if len(got) != 0 {
		t.Errorf("expected empty result for nil input, got %d", len(got))
	}
}
