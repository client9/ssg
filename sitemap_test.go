package ssg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sitemapArtifact(outputFile string) Artifact {
	meta := ContentSourceConfig{}
	if outputFile != "" {
		meta["OutputFile"] = outputFile
	}
	return Artifact{Meta: meta}
}

func TestWriteSitemap_basic(t *testing.T) {
	dir := t.TempDir()
	artifacts := []Artifact{
		sitemapArtifact("index.html"),
		sitemapArtifact("posts/hello/index.html"),
	}

	if err := WriteSitemap(dir, "https://example.com", artifacts); err != nil {
		t.Fatalf("WriteSitemap: %v", err)
	}

	got := readSitemap(t, dir)
	if !strings.Contains(got, "https://example.com/index.html") {
		t.Errorf("missing root page URL:\n%s", got)
	}
	if !strings.Contains(got, "https://example.com/posts/hello/index.html") {
		t.Errorf("missing nested page URL:\n%s", got)
	}
}

func TestWriteSitemap_trailingSlashOnBaseURL(t *testing.T) {
	dir := t.TempDir()
	artifacts := []Artifact{sitemapArtifact("index.html")}

	if err := WriteSitemap(dir, "https://example.com/", artifacts); err != nil {
		t.Fatalf("WriteSitemap: %v", err)
	}

	got := readSitemap(t, dir)
	if strings.Contains(got, "//index.html") {
		t.Errorf("double slash in URL:\n%s", got)
	}
	if !strings.Contains(got, "https://example.com/index.html") {
		t.Errorf("URL malformed:\n%s", got)
	}
}

func TestWriteSitemap_skipsEmptyOutputFile(t *testing.T) {
	dir := t.TempDir()
	artifacts := []Artifact{
		sitemapArtifact("index.html"),
		sitemapArtifact(""), // no OutputFile
		sitemapArtifact("about.html"),
	}

	if err := WriteSitemap(dir, "https://example.com", artifacts); err != nil {
		t.Fatalf("WriteSitemap: %v", err)
	}

	got := readSitemap(t, dir)
	if strings.Count(got, "<loc>") != 2 {
		t.Errorf("expected 2 <loc> entries, got:\n%s", got)
	}
}

func TestWriteSitemap_emptyArtifacts(t *testing.T) {
	dir := t.TempDir()

	if err := WriteSitemap(dir, "https://example.com", nil); err != nil {
		t.Fatalf("WriteSitemap: %v", err)
	}

	got := readSitemap(t, dir)
	if !strings.Contains(got, "<urlset") {
		t.Errorf("expected valid XML even with no artifacts:\n%s", got)
	}
	if strings.Contains(got, "<url>") {
		t.Errorf("expected no <url> entries for empty artifacts:\n%s", got)
	}
}

func TestWriteSitemap_validXML(t *testing.T) {
	dir := t.TempDir()
	artifacts := []Artifact{sitemapArtifact("index.html")}

	if err := WriteSitemap(dir, "https://example.com", artifacts); err != nil {
		t.Fatalf("WriteSitemap: %v", err)
	}

	got := readSitemap(t, dir)
	if !strings.HasPrefix(got, "<?xml") {
		t.Errorf("missing XML declaration:\n%s", got)
	}
	if !strings.Contains(got, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`) {
		t.Errorf("missing sitemap namespace:\n%s", got)
	}
}

func readSitemap(t *testing.T, dir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "sitemap.xml"))
	if err != nil {
		t.Fatalf("reading sitemap.xml: %v", err)
	}
	return string(b)
}
