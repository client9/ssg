package shortcode

import (
	"errors"
	"testing"
)

var errBadInput = errors.New("bad input")

func TestAddErrorAndErr(t *testing.T) {
	ctx := New()
	ctx.Register("fail", func(c *Context, _ []string, _ string) string {
		c.AddError(errBadInput)
		return "[failed]"
	})

	out, err := ctx.RenderDocument("before $fail after")
	if out != "before [failed] after" {
		t.Errorf("unexpected output: %q", out)
	}
	if !errors.Is(err, errBadInput) {
		t.Errorf("expected errBadInput, got %v", err)
	}
}

func TestErrorPosition(t *testing.T) {
	ctx := New()
	ctx.Register("fail", func(c *Context, _ []string, _ string) string {
		c.AddError(errBadInput)
		return ""
	})

	// $fail is at line 2, col 3 ("  $fail")
	_, _ = ctx.RenderDocument("first line\n  $fail\nthird line")

	var pe *PositionError
	if !errors.As(ctx.Err(), &pe) {
		t.Fatalf("expected *PositionError, got %T: %v", ctx.Err(), ctx.Err())
	}
	if pe.Line != 2 || pe.Col != 3 {
		t.Errorf("expected line 2 col 3, got line %d col %d", pe.Line, pe.Col)
	}
	if !errors.Is(pe.Err, errBadInput) {
		t.Errorf("inner error: expected errBadInput, got %v", pe.Err)
	}
}

func TestErrorPositionMultiline(t *testing.T) {
	ctx := New()
	ctx.Register("e1", func(c *Context, _ []string, _ string) string {
		c.AddError(errors.New("first"))
		return ""
	})
	ctx.Register("e2", func(c *Context, _ []string, _ string) string {
		c.AddError(errors.New("second"))
		return ""
	})

	_, _ = ctx.RenderDocument("$e1\nfoo\n$e2")

	errs := ctx.Errs()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}

	pos := []*PositionError{{}, {}}
	for i, e := range errs {
		if !errors.As(e, &pos[i]) {
			t.Fatalf("error %d: expected *PositionError", i)
		}
	}
	if pos[0].Line != 1 || pos[0].Col != 1 {
		t.Errorf("e1: expected line 1 col 1, got line %d col %d", pos[0].Line, pos[0].Col)
	}
	if pos[1].Line != 3 || pos[1].Col != 1 {
		t.Errorf("e2: expected line 3 col 1, got line %d col %d", pos[1].Line, pos[1].Col)
	}
}

func TestMultipleErrors(t *testing.T) {
	ctx := New()
	ctx.Register("e1", func(c *Context, _ []string, _ string) string {
		c.AddError(errors.New("first"))
		return ""
	})
	ctx.Register("e2", func(c *Context, _ []string, _ string) string {
		c.AddError(errors.New("second"))
		return ""
	})

	_, err := ctx.RenderDocument("$e1 $e2")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(ctx.Errs()) != 2 {
		t.Errorf("expected 2 errors, got %d", len(ctx.Errs()))
	}
}

func TestRenderDocumentClearsErrors(t *testing.T) {
	ctx := New()
	ctx.Register("fail", func(c *Context, _ []string, _ string) string {
		c.AddError(errBadInput)
		return ""
	})

	ctx.RenderDocument("$fail")
	if ctx.Err() == nil {
		t.Fatal("expected error after first render")
	}

	// Second call clears previous errors.
	_, err := ctx.RenderDocument("no macros here")
	if err != nil {
		t.Errorf("expected clean slate, got %v", err)
	}
}

func TestNoErrorOnSuccess(t *testing.T) {
	ctx := New()
	ctx.Register("ok", func(_ *Context, _ []string, _ string) string {
		return "fine"
	})

	out, err := ctx.RenderDocument("$ok")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if out != "fine" {
		t.Errorf("unexpected output: %q", out)
	}
}
