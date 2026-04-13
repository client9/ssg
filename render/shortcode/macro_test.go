package shortcode

import (
	"strings"
	"testing"
)

func TestRenderLiteral(t *testing.T) {
	ctx := New()
	got := ctx.Render("hello world")
	if got != "hello world" {
		t.Errorf("got %q", got)
	}
}

func TestDollarEscape(t *testing.T) {
	ctx := New()
	got := ctx.Render("price: $$5.00")
	if got != "price: $5.00" {
		t.Errorf("got %q", got)
	}
}

func TestSimpleCmd(t *testing.T) {
	ctx := New()
	ctx.Register("hr", func(_ *Context, _ []string, _ string) string {
		return "---"
	})
	got := ctx.Render("before $hr after")
	if got != "before --- after" {
		t.Errorf("got %q", got)
	}
}

func TestCmdWithBody(t *testing.T) {
	ctx := New()
	ctx.Register("b", func(_ *Context, _ []string, body string) string {
		return "**" + body + "**"
	})
	got := ctx.Render("hello $b{world}.")
	if got != "hello **world**." {
		t.Errorf("got %q", got)
	}
}

func TestCmdWithArgs(t *testing.T) {
	ctx := New()
	ctx.Register("link", func(_ *Context, args []string, body string) string {
		if len(args) > 0 {
			return "[" + body + "](" + args[0] + ")"
		}
		return body
	})
	got := ctx.Render(`$link[https://example.com]{click here}`)
	want := "[click here](https://example.com)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCmdWithQuotedArgs(t *testing.T) {
	ctx := New()
	ctx.Register("tag", func(_ *Context, args []string, _ string) string {
		return strings.Join(args, "|")
	})
	got := ctx.Render(`$tag["first" "second" "third"]`)
	want := "first|second|third"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCmdWithNamedArgs(t *testing.T) {
	ctx := New()
	ctx.Register("img", func(_ *Context, args []string, _ string) string {
		attrs := ParseNamedArgs(args)
		return `<img src="` + attrs["src"] + `" alt="` + attrs["alt"] + `">`
	})
	got := ctx.Render(`$img[src=photo.jpg alt="a photo"]`)
	want := `<img src="photo.jpg" alt="a photo">`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRecursiveRender(t *testing.T) {
	ctx := New()
	ctx.Register("h1", func(c *Context, _ []string, body string) string {
		return "<h1>" + c.Render(body) + "</h1>"
	})
	ctx.Register("em", func(c *Context, _ []string, body string) string {
		return "<em>" + c.Render(body) + "</em>"
	})
	got := ctx.Render("Before\n$h1{this is $em{a} headline}\nAfter")
	want := "Before\n<h1>this is <em>a</em> headline</h1>\nAfter"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNestedBraces(t *testing.T) {
	ctx := New()
	ctx.Register("raw", func(_ *Context, _ []string, body string) string {
		return body
	})
	got := ctx.Render(`$raw{a{b}c}`)
	if got != "a{b}c" {
		t.Errorf("got %q", got)
	}
}

func TestUnknownCmdPassthrough(t *testing.T) {
	ctx := New()
	got := ctx.Render("$unknown{text}")
	if got != "$unknown{text}" {
		t.Errorf("got %q", got)
	}
}

func TestHeadlineExample(t *testing.T) {
	ctx := New()
	ctx.Register("h1", func(c *Context, _ []string, body string) string {
		return "<h1>" + c.Render(body) + "</h1>"
	})
	input := "Before\n$h1{this is a headline}\nAfter"
	want := "Before\n<h1>this is a headline</h1>\nAfter"
	got := ctx.Render(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
