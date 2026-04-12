package ssg

import (
	"testing"
)

func TestFilterPages_keepAll(t *testing.T) {
	pages := []ContentSourceConfig{
		{"Title": "A"},
		{"Title": "B"},
	}
	got := FilterPages(pages, func(ContentSourceConfig) bool { return true })
	if len(got) != 2 {
		t.Errorf("expected 2 pages, got %d", len(got))
	}
}

func TestFilterPages_removeAll(t *testing.T) {
	pages := []ContentSourceConfig{
		{"Title": "A"},
		{"Title": "B"},
	}
	got := FilterPages(pages, func(ContentSourceConfig) bool { return false })
	if len(got) != 0 {
		t.Errorf("expected 0 pages, got %d", len(got))
	}
}

func TestFilterPages_draftExclusion(t *testing.T) {
	pages := []ContentSourceConfig{
		{"Title": "Published"},
		{"Title": "Draft", "draft": true},
		{"Title": "Also published", "draft": false},
	}
	got := FilterPages(pages, func(p ContentSourceConfig) bool {
		draft, _ := p["draft"].(bool)
		return !draft
	})
	if len(got) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(got))
	}
	for _, p := range got {
		if p["Title"] == "Draft" {
			t.Error("draft page should have been excluded")
		}
	}
}

func TestFilterPages_byStringField(t *testing.T) {
	pages := []ContentSourceConfig{
		{"section": "blog"},
		{"section": "docs"},
		{"section": "blog"},
	}
	got := FilterPages(pages, func(p ContentSourceConfig) bool {
		return p.Get("section") == "blog"
	})
	if len(got) != 2 {
		t.Errorf("expected 2 blog pages, got %d", len(got))
	}
}

func TestFilterPages_empty(t *testing.T) {
	got := FilterPages(nil, func(ContentSourceConfig) bool { return true })
	if len(got) != 0 {
		t.Errorf("expected empty result for nil input, got %d", len(got))
	}
}

func TestFilterPages_doesNotModifyInput(t *testing.T) {
	pages := []ContentSourceConfig{
		{"Title": "A"},
		{"Title": "B"},
		{"Title": "C"},
	}
	original := len(pages)
	FilterPages(pages, func(ContentSourceConfig) bool { return false })
	if len(pages) != original {
		t.Errorf("FilterPages modified the input slice")
	}
}
