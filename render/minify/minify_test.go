package minify

import (
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
	cfg := ssg.ContentSourceConfig{"OutputFile": "index.html"}
	out, err := r.Run(nil, cfg, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.([]byte)
	if len(got) >= len(in) {
		t.Errorf("expected minified output to be shorter; got %q", got)
	}
}

func TestMinifyPassthrough(t *testing.T) {
	r := New()
	in := []byte("hello world")
	cfg := ssg.ContentSourceConfig{"OutputFile": "file.unknown"}
	out, err := r.Run(nil, cfg, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.([]byte)
	if string(got) != string(in) {
		t.Errorf("expected passthrough; got %q", got)
	}
}
