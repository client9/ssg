package shortcode

import "testing"

func imgNamed(_ *Context, named map[string]string, _ string) string {
	return `<img src="` + named["src"] + `" alt="` + named["alt"] + `">`
}

// Named args by key=value.
func TestMakeNamedExplicit(t *testing.T) {
	ctx := New()
	ctx.RegisterNamed("img", imgNamed, "src", "alt")

	got := ctx.Render(`$img[src=photo.jpg alt="a photo"]`)
	want := `<img src="photo.jpg" alt="a photo">`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Positional args mapped to names by position.
func TestMakeNamedPositional(t *testing.T) {
	ctx := New()
	ctx.RegisterNamed("img", imgNamed, "src", "alt")

	got := ctx.Render(`$img[photo.jpg "a photo"]`)
	want := `<img src="photo.jpg" alt="a photo">`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Mixed: first arg positional, second named.
func TestMakeNamedMixed(t *testing.T) {
	ctx := New()
	ctx.RegisterNamed("img", imgNamed, "src", "alt")

	got := ctx.Render(`$img[photo.jpg alt="a photo"]`)
	want := `<img src="photo.jpg" alt="a photo">`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Extra positional args beyond paramNames are dropped.
func TestMakeNamedExtraArgs(t *testing.T) {
	ctx := New()
	ctx.RegisterNamed("img", imgNamed, "src") // only "src" named

	got := ctx.Render(`$img[photo.jpg extra ignored]`)
	want := `<img src="photo.jpg" alt="">`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// No paramNames: only explicit key=value args work.
func TestMakeNamedNoParamNames(t *testing.T) {
	ctx := New()
	ctx.RegisterNamed("img", imgNamed) // no positional mapping

	got := ctx.Render(`$img[src=photo.jpg alt="a photo"]`)
	want := `<img src="photo.jpg" alt="a photo">`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
