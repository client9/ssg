package main

import (
	"fmt"
	"github.com/client9/ssg"
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

	// create page assembly templates
	pageTemplate, err := ssg.NewPageRender("layout", fns)
	if err != nil {
		log.Fatalf("Page Template failed: %v", err)
	}

	// config and pipeline
	conf := ssg.SiteConfig{
		OutputDir: "public",
		Pipeline: []ssg.Renderer{
			ssg.NewTemplateMacro(fns),
			ssg.HTMLRender,
			pageTemplate,
			HTMLPretty,
			ssg.WriteOutput("public"),
		},
	}

	//  create array of pages
	//  One may manually create various pages
	//  from database or something else

	pages := []ssg.ContentSourceConfig{}

	// do it
	err = ssg.Main2(conf, &pages)
	if err != nil {
		log.Fatalf("Main failed: %s", err)
	}
}
