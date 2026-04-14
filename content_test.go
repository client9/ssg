package ssg

import (
	"os"
	"path/filepath"
	"testing"
)

// writeContent creates a temporary content directory with the given files.
func writeContent(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for rel, body := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0666); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestFileWalker_basic(t *testing.T) {
	dir := writeContent(t, map[string]string{
		"post.html": "body only",
	})

	rules := []Rule{
		{Pattern: "**/*.html", Loader: Passthrough, Pipeline: Pipeline{name: "nothing", stages: nil}},
	}

	var artifacts []Artifact
	if err := FileWalker(dir, rules)(nil, &artifacts); err != nil {
		t.Fatalf("FileWalker: %v", err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("expected 1 artifact, got %d", len(artifacts))
	}
	if string(artifacts[0].Meta["Content"].([]byte)) != "body only" {
		t.Errorf("unexpected content: %s", artifacts[0].Meta["Content"])
	}
}

func TestFileWalker_noMatch(t *testing.T) {
	dir := writeContent(t, map[string]string{
		"post.md": "# hello",
	})

	rules := []Rule{
		{Pattern: "**/*.html", Loader: Passthrough, Pipeline: Pipeline{name: "nothing", stages: nil}},
	}

	var artifacts []Artifact
	if err := FileWalker(dir, rules)(nil, &artifacts); err != nil {
		t.Fatalf("FileWalker: %v", err)
	}
	if len(artifacts) != 0 {
		t.Errorf("expected 0 artifacts, got %d", len(artifacts))
	}
}

func TestFileWalker_nilLoader(t *testing.T) {
	dir := writeContent(t, map[string]string{
		"_draft.html": "skip me",
	})

	rules := []Rule{
		{Pattern: "**/_*"}, // nil Loader = skip
	}

	var artifacts []Artifact
	if err := FileWalker(dir, rules)(nil, &artifacts); err != nil {
		t.Fatalf("FileWalker: %v", err)
	}
	if len(artifacts) != 0 {
		t.Errorf("expected 0 artifacts for nil loader, got %d", len(artifacts))
	}
}

func TestFileWalker_skipHiddenDirs(t *testing.T) {
	dir := writeContent(t, map[string]string{
		".git/config": "gitconfig",
		"post.html":   "visible",
	})

	rules := []Rule{
		{Pattern: "**/*.html", Loader: Passthrough, Pipeline: Pipeline{name: "nothing", stages: nil}},
	}

	var artifacts []Artifact
	if err := FileWalker(dir, rules)(nil, &artifacts); err != nil {
		t.Fatalf("FileWalker: %v", err)
	}
	if len(artifacts) != 1 {
		t.Errorf("expected 1 artifact (hidden dir skipped), got %d", len(artifacts))
	}
}

func TestFileWalker_independentMeta(t *testing.T) {
	dir := writeContent(t, map[string]string{
		"a.md": "# a",
		"b.md": "# b",
	})

	rules := []Rule{
		{Pattern: "**/*.md", Loader: Passthrough, Pipeline: Pipeline{name: "nothing", stages: nil}},
	}

	var artifacts []Artifact
	if err := FileWalker(dir, rules)(nil, &artifacts); err != nil {
		t.Fatalf("FileWalker: %v", err)
	}
	if len(artifacts) != 2 {
		t.Fatalf("expected 2 artifacts (one per file), got %d", len(artifacts))
	}
	// Each artifact must have an independent meta map.
	artifacts[0].Meta["OutputFile"] = "a.html"
	artifacts[1].Meta["OutputFile"] = "b.html"
	if artifacts[0].Meta.OutputFile() == artifacts[1].Meta.OutputFile() {
		t.Error("artifacts should have independent meta maps")
	}
}

func TestFileWalker_firstRuleWins(t *testing.T) {
	dir := writeContent(t, map[string]string{
		"post.html": "content",
	})

	called := 0
	countLoader := MetaLoader(func(raw []byte) (map[string]any, []byte, error) {
		called++
		return map[string]any{}, raw, nil
	})

	rules := []Rule{
		{Pattern: "**/*.html", Loader: countLoader, Pipeline: Pipeline{name: "nothing", stages: nil}},
		{Pattern: "**/*.html", Loader: countLoader, Pipeline: Pipeline{name: "nothing", stages: nil}},
	}

	var artifacts []Artifact
	if err := FileWalker(dir, rules)(nil, &artifacts); err != nil {
		t.Fatalf("FileWalker: %v", err)
	}
	if called != 1 {
		t.Errorf("expected loader called once (first match wins), got %d", called)
	}
}
