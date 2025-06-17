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
type Renderer func(wr io.Writer, source []byte, data any) error

func MultiRender(renders []Renderer, src []byte, data any) ([]byte, error) {
	dest := bytes.Buffer{}
	dest.Grow(len(src) * 2)
	for _, r := range renders {
		if err := r(&dest, src, data); err != nil {
			return nil, err
		}
		src = bytes.Clone(dest.Bytes())
		dest.Reset()
	}
	return src, nil
}

// create a simple macro maker.  You pass-in whatever funcs.
//
//	each page is parsed as a text/template then executed
func NewTemplateMacro(funcs template.FuncMap) Renderer {
	t := template.New("_macro")
	if funcs != nil {
		t = t.Funcs(funcs)
	}
	return func(wr io.Writer, source []byte, data any) error {
		t, err := t.Parse(string(source))
		if err != nil {
			return err
		}
		return t.Execute(wr, data)
	}
}

func Identity(wr io.Writer, source []byte, data any) error {
	_, err := wr.Write(source)
	return err
}

// This parses in the context html and then writes it back out.
// This normalized HTML and prevents the content from breaking the page layout.
func HTMLRender(wr io.Writer, source []byte, data any) error {
	// should this actually be parse fragement?
	// create a div node?
	div := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	}
	children, err := html.ParseFragment(bytes.NewReader(source), div)
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
