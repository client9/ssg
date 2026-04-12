package ssg

import (
	"io"
	"testing"
)

// noopRenderer is a Renderer that discards all input, used to exercise Render
// without needing a real pipeline.
func noopRenderer(wr io.Writer, src io.Reader, data any) error {
	_, err := io.Copy(io.Discard, src)
	return err
}

// minimalPage returns a ContentSourceConfig with the minimum keys required by
// Render (Content and OutputFile).
func minimalPage(title string) ContentSourceConfig {
	return ContentSourceConfig{
		"Content":    []byte("body"),
		"OutputFile": "out.html",
		"InputFile":  "content/page.html",
		"Title":      title,
	}
}

func TestRender_nilGlobals(t *testing.T) {
	pipeline := []Renderer{noopRenderer}
	pages := []ContentSourceConfig{minimalPage("A"), minimalPage("B")}

	if err := Render(pipeline, pages, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRender_globalsInjected(t *testing.T) {
	pipeline := []Renderer{noopRenderer}
	pages := []ContentSourceConfig{minimalPage("A"), minimalPage("B")}
	globals := map[string]any{"Nav": "top-nav"}

	if err := Render(pipeline, pages, globals); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, p := range pages {
		if got := p.Get("Nav"); got != "top-nav" {
			t.Errorf("pages[%d]: expected Nav=%q, got %q", i, "top-nav", got)
		}
	}
}

func TestRender_pageFrontmatterWins(t *testing.T) {
	pipeline := []Renderer{noopRenderer}
	page := minimalPage("A")
	page["Nav"] = "page-nav"
	pages := []ContentSourceConfig{page}
	globals := map[string]any{"Nav": "global-nav"}

	if err := Render(pipeline, pages, globals); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := pages[0].Get("Nav"); got != "page-nav" {
		t.Errorf("expected page frontmatter to win: got %q, want %q", got, "page-nav")
	}
}

func TestRender_globalsAvailableToAllPages(t *testing.T) {
	pipeline := []Renderer{noopRenderer}
	pages := []ContentSourceConfig{minimalPage("A"), minimalPage("B"), minimalPage("C")}
	globals := map[string]any{"Shared": "yes"}

	if err := Render(pipeline, pages, globals); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, p := range pages {
		if got := p.Get("Shared"); got != "yes" {
			t.Errorf("pages[%d]: expected Shared=%q, got %q", i, "yes", got)
		}
	}
}

func TestRender_emptyGlobals(t *testing.T) {
	pipeline := []Renderer{noopRenderer}
	pages := []ContentSourceConfig{minimalPage("A")}

	if err := Render(pipeline, pages, map[string]any{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
