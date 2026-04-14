package ssg

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

// makeLayout creates a temporary layout directory tree for testing.
// files is a map of relative path → template content.
func makeLayout(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0666); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestExecuteTemplate_rootTemplate(t *testing.T) {
	root := makeLayout(t, map[string]string{
		"base.html": `BASE:{{.Content}}`,
	})

	router, err := templateMap(root, nil)
	if err != nil {
		t.Fatalf("templateMap: %v", err)
	}

	var out bytes.Buffer
	data := ContentSourceConfig{"Content": "hello"}
	if err := router.ExecuteTemplate(&out, "base.html", data); err != nil {
		t.Fatalf("ExecuteTemplate: %v", err)
	}
	if got, want := out.String(), "BASE:hello"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExecuteTemplate_subdirectory(t *testing.T) {
	root := makeLayout(t, map[string]string{
		"base.html":        `BASE:{{.Content}}`,
		"blog/single.html": `BLOG:{{.Title}}`,
	})

	router, err := templateMap(root, nil)
	if err != nil {
		t.Fatalf("templateMap: %v", err)
	}

	var out bytes.Buffer
	data := ContentSourceConfig{"Title": "My Post", "Content": "body"}
	if err := router.ExecuteTemplate(&out, "blog/single.html", data); err != nil {
		t.Fatalf("ExecuteTemplate: %v", err)
	}
	if got, want := out.String(), "BLOG:My Post"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestExecuteTemplate_inheritance verifies that a child template can call a
// template defined in a parent directory.
func TestExecuteTemplate_inheritance(t *testing.T) {
	root := makeLayout(t, map[string]string{
		// Parent defines the outer wrapper.
		"base.html": `<html>{{.Content}}</html>`,
		// Child wraps its own content using the parent template.
		"blog/single.html": `{{template "base.html" .}}`,
	})

	router, err := templateMap(root, nil)
	if err != nil {
		t.Fatalf("templateMap: %v", err)
	}

	var out bytes.Buffer
	data := ContentSourceConfig{"Content": "article body"}
	if err := router.ExecuteTemplate(&out, "blog/single.html", data); err != nil {
		t.Fatalf("ExecuteTemplate: %v", err)
	}
	if got, want := out.String(), "<html>article body</html>"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestExecuteTemplate_deepInheritance verifies a two-level directory can call
// templates from both its parent and grandparent.
func TestExecuteTemplate_deepInheritance(t *testing.T) {
	root := makeLayout(t, map[string]string{
		"wrap.html":            `[{{.Content}}]`,
		"blog/mid.html":        `MID:{{.Content}}`,
		"blog/posts/page.html": `{{template "wrap.html" .}}`,
	})

	router, err := templateMap(root, nil)
	if err != nil {
		t.Fatalf("templateMap: %v", err)
	}

	var out bytes.Buffer
	data := ContentSourceConfig{"Content": "deep"}
	if err := router.ExecuteTemplate(&out, "blog/posts/page.html", data); err != nil {
		t.Fatalf("ExecuteTemplate: %v", err)
	}
	if got, want := out.String(), "[deep]"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExecuteTemplate_unknownDirectory(t *testing.T) {
	root := makeLayout(t, map[string]string{
		"base.html": `hello`,
	})

	router, err := templateMap(root, nil)
	if err != nil {
		t.Fatalf("templateMap: %v", err)
	}

	var out bytes.Buffer
	err = router.ExecuteTemplate(&out, "nope/page.html", nil)
	if err == nil {
		t.Fatal("expected error for unknown directory, got nil")
	}
	if !strings.Contains(err.Error(), "nope") {
		t.Errorf("error should mention the missing directory, got: %v", err)
	}
}

func TestExecuteTemplate_customFuncs(t *testing.T) {
	root := makeLayout(t, map[string]string{
		"base.html": `{{shout .Content}}`,
	})

	fns := template.FuncMap{
		"shout": func(s string) string { return strings.ToUpper(s) + "!" },
	}
	router, err := templateMap(root, fns)
	if err != nil {
		t.Fatalf("templateMap: %v", err)
	}

	var out bytes.Buffer
	data := ContentSourceConfig{"Content": "hello"}
	if err := router.ExecuteTemplate(&out, "base.html", data); err != nil {
		t.Fatalf("ExecuteTemplate: %v", err)
	}
	if got, want := out.String(), "HELLO!"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNewPageRender(t *testing.T) {
	root := makeLayout(t, map[string]string{
		"page.html": `<title>{{.Title}}</title><body>{{.Content}}</body>`,
	})

	stage, err := NewPageRender(root, nil)
	if err != nil {
		t.Fatalf("NewPageRender: %v", err)
	}

	cfg := ContentSourceConfig{
		"TemplateName": "page.html",
		"Title":        "Hello",
	}
	p := Pipeline{
		name:   "test",
		stages: []Stage{stage},
	}

	got, err := RunPipeline[[]byte](nil, cfg, p, []byte("<p>world</p>"))
	if err != nil {
		t.Fatalf("EvalStages: %v", err)
	}

	want := "<title>Hello</title><body><p>world</p></body>"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
