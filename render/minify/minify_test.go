package minify

import (
	"bytes"
	"testing"

	"github.com/client9/ssg"
)

func TestMimeFromFile(t *testing.T) {
	cases := []struct {
		file string
		want string
	}{
		{"index.html", "text/html"},
		{"style.css", "text/css"},
		{"app.js", "text/javascript"},
		{"logo.svg", "image/svg+xml"},
		{"data.json", "application/json"},
		{"feed.xml", "application/xml"},
		{"unknown.xyz", ""},
		{"app.mjs", "text/javascript"},
		{"noextension", ""},
	}
	for _, c := range cases {
		got := mimeFromFile(c.file)
		if got != c.want {
			t.Errorf("mimeFromFile(%q) = %q, want %q", c.file, got, c.want)
		}
	}
}

func TestMinifyHTML(t *testing.T) {
	r := New()
	in := []byte("  <html>  <body>  <p>  hello  </p>  </body>  </html>  ")
	out := &bytes.Buffer{}
	data := ssg.ContentSourceConfig{"OutputFile": "index.html"}
	if err := r(out, bytes.NewReader(in), data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.String()
	if len(got) >= len(in) {
		t.Errorf("expected minified output to be shorter; got %q", got)
	}
}

func TestMinifyPassthrough(t *testing.T) {
	r := New()
	in := []byte("hello world")
	out := &bytes.Buffer{}
	data := ssg.ContentSourceConfig{"OutputFile": "file.unknown"}
	if err := r(out, bytes.NewReader(in), data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(out.Bytes(), in) {
		t.Errorf("expected passthrough; got %q", out.String())
	}
}
