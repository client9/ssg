package markdowncontent

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"io"

	"github.com/client9/ssg"
)

func New() ssg.Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	return NewGoldmark(md)
}

func NewGoldmark(md goldmark.Markdown) ssg.Renderer {
	return func(wr io.Writer, src io.Reader, data any) error {
		buf := bytes.Buffer{}
		_, err := io.Copy(&buf, src)
		if err != nil {
			return err
		}
		return md.Convert(buf.Bytes(), wr)
	}
}
