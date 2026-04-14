package ssg

import (
	"fmt"
	"testing"
)

func TestStep_typeMatch(t *testing.T) {
	s := Step("double", func(_ *Context, _ ContentSourceConfig, in int) (int, error) {
		return in * 2, nil
	})
	out, err := s.Run(nil, nil, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.(int) != 6 {
		t.Errorf("got %v, want 6", out)
	}
}

func TestStep_typeMismatch(t *testing.T) {
	s := Step("wants-int", func(_ *Context, _ ContentSourceConfig, in int) (int, error) {
		return in, nil
	})
	_, err := s.Run(nil, nil, "not an int")
	if err == nil {
		t.Fatal("expected type mismatch error, got nil")
	}
}

func TestEvalStages_chain(t *testing.T) {
	p := Pipeline{
		name: "test",
		stages: []Stage{
			Step("to-string", func(_ *Context, _ ContentSourceConfig, in int) (string, error) {
				return fmt.Sprintf("%d", in), nil
			}),
			Step("append", func(_ *Context, _ ContentSourceConfig, in string) (string, error) {
				return in + "!", nil
			}),
		},
	}
	got, err := RunPipeline[string](nil, nil, p, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "42!" {
		t.Errorf("got %q, want %q", got, "42!")
	}
}

func TestEvalStages_stageError(t *testing.T) {
	p := Pipeline{
		name: "test",
		stages: []Stage{
			Step("fail", func(_ *Context, _ ContentSourceConfig, in string) (string, error) {
				return "", fmt.Errorf("boom")
			}),
		},
	}
	_, err := RunPipeline[string](nil, nil, p, "x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := `pipeline "test", stage "fail" failed: boom`
	if err.Error() != want {
		t.Errorf("unexpected error format: got %q, want %q", err.Error(), want)
	}
}

func noopStage() Stage {
	return Step("noop", func(_ *Context, _ ContentSourceConfig, in []byte) ([]byte, error) {
		return in, nil
	})
}

func minimalArtifact(title string) Artifact {
	return Artifact{
		Meta: ContentSourceConfig{
			"Content":    []byte("body"),
			"OutputFile": "out.html",
			"InputFile":  "content/page.html",
			"Title":      title,
		},
		Pipeline: Pipeline{
			name:   "test",
			stages: []Stage{noopStage()},
		},
	}
}

func TestRender_basic(t *testing.T) {
	artifacts := []Artifact{minimalArtifact("A"), minimalArtifact("B")}
	if err := Render(nil, &artifacts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRender_globalsInjected(t *testing.T) {
	artifacts := []Artifact{minimalArtifact("A"), minimalArtifact("B")}
	ctx := &Context{Globals: map[string]any{"Nav": "top-nav"}}

	if err := Render(ctx, &artifacts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, a := range artifacts {
		if got := a.Meta.Get("Nav"); got != "top-nav" {
			t.Errorf("artifacts[%d]: expected Nav=%q, got %q", i, "top-nav", got)
		}
	}
}

func TestRender_pageFrontmatterWins(t *testing.T) {
	a := minimalArtifact("A")
	a.Meta["Nav"] = "page-nav"
	artifacts := []Artifact{a}
	ctx := &Context{Globals: map[string]any{"Nav": "global-nav"}}

	if err := Render(ctx, &artifacts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := artifacts[0].Meta.Get("Nav"); got != "page-nav" {
		t.Errorf("expected page frontmatter to win: got %q, want %q", got, "page-nav")
	}
}
