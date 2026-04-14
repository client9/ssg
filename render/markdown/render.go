package markdown

import (
	"bytes"

	"github.com/client9/ssg"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// New returns a DynStage that converts Markdown to HTML using Goldmark
// with GitHub Flavored Markdown extensions and auto heading IDs.
func New() ssg.Stage {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	return NewGoldmark(md)
}

// NewGoldmark wraps a custom Goldmark instance as a DynStage.
func NewGoldmark(md goldmark.Markdown) ssg.Stage {
	return ssg.Step("markdown", func(_ *ssg.Context, _ ssg.ContentSourceConfig, in []byte) ([]byte, error) {
		var buf bytes.Buffer
		if err := md.Convert(in, &buf); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	})
}
