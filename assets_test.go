package ssg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyAssets(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	// Create source tree
	write := func(rel, content string) {
		t.Helper()
		full := filepath.Join(src, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0666); err != nil {
			t.Fatal(err)
		}
	}

	write("style.css", "body{}")
	write("img/logo.svg", "<svg/>")
	write(".hidden", "skip me")
	write(".dotdir/file.css", "skip me too")

	if err := CopyAssets(src, dest); err != nil {
		t.Fatalf("CopyAssets error: %v", err)
	}

	check := func(rel, want string) {
		t.Helper()
		got, err := os.ReadFile(filepath.Join(dest, rel))
		if err != nil {
			t.Errorf("missing %s: %v", rel, err)
			return
		}
		if string(got) != want {
			t.Errorf("%s: got %q, want %q", rel, got, want)
		}
	}
	absent := func(rel string) {
		t.Helper()
		if _, err := os.Stat(filepath.Join(dest, rel)); err == nil {
			t.Errorf("expected %s to be absent", rel)
		}
	}

	check("style.css", "body{}")
	check("img/logo.svg", "<svg/>")
	absent(".hidden")
	absent(".dotdir/file.css")
}
