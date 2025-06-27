package ssg

import (
	"bytes"
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
