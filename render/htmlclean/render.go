package htmlclean

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Render parses src as an HTML fragment, normalizes it, and writes it to wr.
// This prevents malformed content from breaking the surrounding page layout.
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
		if err = html.Render(wr, n); err != nil {
			return err
		}
	}
	return nil
}
