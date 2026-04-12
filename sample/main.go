package main

import (
	"fmt"
	"github.com/client9/ssg"
	"github.com/client9/ssg/render/htmlclean"
	"html"
	"io"
	"log"
	"strings"
	"text/template"

	"github.com/yosssi/gohtml"
)

// same "content macro" using golang templates
//
// {{ elink "https://.../"  "name" }}
// creates a <a> that opens in a new window
func elink(href string, body string) string {
	return fmt.Sprintf("<a href=%q target=_blank>%s</a>",
		html.EscapeString(href),
		html.EscapeString(body))
}

// here's an example of a post processor
func HTMLPretty(wr io.Writer, source io.Reader, data any) error {
	src, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	wr.Write(gohtml.FormatBytes(src))
	return nil
}

func main() {
	// various golang template functions
	fns := template.FuncMap{
		"upper": strings.ToUpper,
		"elink": elink,
	}

	// file loading config
	loadConf := ssg.LoadConfig{
		ContentDir:   "content",
		BaseTemplate: "baseof.html",
		MetaSplit:    ssg.MetaSplitJson,
		MetaParser:   ssg.MetaParseJson,
		InputExt:     ".html",
		OutputExt:    ".html",
		IndexSource:  "index.html",
		IndexDest:    "index.html",
	}

	// rendering pipeline
	pipeline := []ssg.Renderer{
		ssg.NewTemplateMacro(fns),
		htmlclean.Render,
		ssg.Must(ssg.NewPageRender("layout", fns)),
		HTMLPretty,
		ssg.WriteOutput("public"),
	}

	pages := []ssg.ContentSourceConfig{}

	if err := ssg.LoadContent(loadConf, &pages); err != nil {
		log.Fatalf("load content failed: %s", err)
	}

	if err := ssg.Render(pipeline, pages, nil); err != nil {
		log.Fatalf("Main failed: %s", err)
	}
}
