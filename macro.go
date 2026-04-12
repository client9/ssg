package ssg

import (
	"io"
	"text/template"
)

// NewTemplateMacro returns a Renderer that treats each page's content as a
// text/template, executing it with the page data. This is intended for
// lightweight macros embedded in markdown or plain-text content — for example,
// a shortcode-style helper that expands to an HTML snippet.
//
// Use text/template (not html/template) because the author controls the input
// and the output is typically raw markup, not user-supplied data.
func NewTemplateMacro(funcs template.FuncMap) Renderer {
	t := template.New("_macro")
	if funcs != nil {
		t = t.Funcs(funcs)
	}
	return func(wr io.Writer, src io.Reader, data any) error {
		raw, err := io.ReadAll(src)
		if err != nil {
			return err
		}
		t, err = t.Parse(string(raw))
		if err != nil {
			return err
		}
		return t.Execute(wr, data)
	}
}
