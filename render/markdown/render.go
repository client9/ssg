package markdown

import (
	"bytes"
	"io"

	"github.com/client9/ssg"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// New returns a Renderer that converts Markdown to HTML using Goldmark
// with GitHub Flavored Markdown extensions and auto heading IDs.
func New() ssg.Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	return NewGoldmark(md)
}

// NewGoldmark wraps a custom Goldmark instance as a Renderer.
func NewGoldmark(md goldmark.Markdown) ssg.Renderer {
	return func(wr io.Writer, src io.Reader, data any) error {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, src); err != nil {
			return err
		}
		return md.Convert(buf.Bytes(), wr)
	}
}
