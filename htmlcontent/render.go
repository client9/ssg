package htmlcontent

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
)

// This parses in the context html and then writes it back out.
// This normalized HTML and prevents the content from breaking the page layout.
//
// This is more a sample of the Render process
//
func Render(wr io.Writer, src io.Reader, data any) error {
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
