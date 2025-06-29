package main

import (
	"fmt"
	"github.com/client9/ssg"
	"github.com/client9/ssg/htmlcontent"
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

	// config and pipeline
	conf := ssg.SiteConfig{
		ContentDir: "content",
		BaseTemplate: "baseof.html",
		MetaSplit: ssg.MetaSplitJson,
		MetaParser: ssg.MetaParseJson,
		InputExt: ".html",
		OutputExt: ".html",
		IndexSource: "index.html",
		IndexDest: "index.html",
		Pipeline: []ssg.Renderer{
			ssg.NewTemplateMacro(fns),
			htmlcontent.Render,
			ssg.Must(ssg.NewPageRender("layout", fns)),
			HTMLPretty,
			ssg.WriteOutput("public"),
		},
	}

	//  create array of pages
	//  One may manually create various pages
	//  from database or something else

	pages := []ssg.ContentSourceConfig{}
 
	// load in content
	if err := ssg.LoadContent(conf, &pages); err != nil {
		log.Fatalf("load content failed: %s", err)
	}

	// do it
	if err := ssg.Main2(conf, &pages); err != nil {
		log.Fatalf("Main failed: %s", err)
	}
}
