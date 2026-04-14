package htmlclean

import (
	"bytes"

	"github.com/client9/ssg"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Render is a DynStage that parses input as an HTML fragment, normalizes it,
// and returns the result. This prevents malformed content from breaking the
// surrounding page layout.
var Render = ssg.Step("htmlclean", func(_ *ssg.Context, _ ssg.ContentSourceConfig, in []byte) ([]byte, error) {
	div := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	}
	children, err := html.ParseFragment(bytes.NewReader(in), div)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	for _, n := range children {
		if err = html.Render(&buf, n); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
})
