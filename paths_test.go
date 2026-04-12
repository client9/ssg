package ssg

import "testing"

func TestCleanURLs(t *testing.T) {
	tr := CleanURLs(".md", ".html")
	tests := []struct {
		in   string
		want string
	}{
		{"foo.md", "foo/index.html"},
		{"bar/baz.md", "bar/baz/index.html"},
		{"index.md", "index.html"},           // root index stays flat
		{"blog/index.md", "blog/index.html"}, // nested index stays flat
		{"foo.html", ""},                     // wrong extension → skip
		{"foo.txt", ""},                      // wrong extension → skip
	}
	for _, tt := range tests {
		got := tr(tt.in)
		if got != tt.want {
			t.Errorf("CleanURLs(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestUglyURLs(t *testing.T) {
	tr := UglyURLs(".md", ".html")
	tests := []struct {
		in   string
		want string
	}{
		{"foo.md", "foo.html"},
		{"bar/baz.md", "bar/baz.html"},
		{"index.md", "index.html"},
		{"foo.txt", ""}, // wrong extension → skip
	}
	for _, tt := range tests {
		got := tr(tt.in)
		if got != tt.want {
			t.Errorf("UglyURLs(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSlugNormalize(t *testing.T) {
	tr := SlugNormalize(UglyURLs(".md", ".html"))
	tests := []struct {
		in   string
		want string
	}{
		{"foo.md", "foo.html"},
		{"Foo Bar.md", "foo-bar.html"},
		{"my_post.md", "my-post.html"},
		{"My_Post.md", "my-post.html"},
		{"blog/Foo Bar.md", "blog/foo-bar.html"},
	}
	for _, tt := range tests {
		got := tr(tt.in)
		if got != tt.want {
			t.Errorf("SlugNormalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSlugNormalize_withCleanURLs(t *testing.T) {
	tr := SlugNormalize(CleanURLs(".md", ".html"))
	if got := tr("My Post.md"); got != "my-post/index.html" {
		t.Errorf("got %q, want %q", got, "my-post/index.html")
	}
}
