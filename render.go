package ssg

import (
	"bytes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"text/template"
)

// "Template Macro"  processes source and writing output
//
//	"data" is optional supplemental data.  It is just passed on to the
//
// underlying implementation.  If not used, then data can be nil.
type Renderer func(wr io.Writer, src io.Reader, data any) error

func MultiRender(renders []Renderer, initial []byte, data any) error {
	src := bytes.NewBuffer(initial)
	src.Grow(len(initial) * 2)
	dest := &bytes.Buffer{}
	dest.Grow(src.Cap())

	for _, r := range renders {
		if err := r(dest, src, data); err != nil {
			return err
		}
		src, dest = dest, src
		dest.Reset()
	}
	return nil
}

// create a simple macro maker.  You pass-in whatever funcs.
//
//	each page is parsed as a text/template then executed
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

func ToBytes(buf *bytes.Buffer) Renderer {
	return func(wr io.Writer, src io.Reader, data any) error {
		_, err := io.Copy(buf, src)
		return err
	}
}

func Identity(wr io.Writer, src io.Reader, data any) error {
	_, err := io.Copy(wr, src)
	return err
}

// This parses in the context html and then writes it back out.
// This normalized HTML and prevents the content from breaking the page layout.
func HTMLRender(wr io.Writer, src io.Reader, data any) error {
	// should this actually be parse fragement?
	// create a div node?
	div := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	}
	children, err := html.ParseFragment(src, div)
	if err != nil {
		return err
	}
	for _, n := range children {
		err = html.Render(wr, n)
		if err != nil {
			return err
		}
	}
	return nil
}
