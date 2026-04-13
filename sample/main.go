package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"sort"
	"strings"
	"text/template"

	"github.com/client9/ssg"
	metajson "github.com/client9/ssg/meta/json"
	"github.com/client9/ssg/render/htmlclean"
	"github.com/yosssi/gohtml"
)

// {{ elink "https://.../" "name" }} creates an <a> that opens in a new window.
func elink(href string, body string) string {
	return fmt.Sprintf("<a href=%q target=_blank>%s</a>",
		html.EscapeString(href),
		html.EscapeString(body))
}

func HTMLPretty(wr io.Writer, source io.Reader, data any) error {
	src, err := io.ReadAll(source)
	if err != nil {
		return err
	}
	wr.Write(gohtml.FormatBytes(src))
	return nil
}

// slug converts a tag name to a URL-safe string.
func slug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func tmplMacro(funcs template.FuncMap) ssg.Renderer {
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

func main() {
	fns := template.FuncMap{
		"upper": strings.ToUpper,
		"elink": elink,
	}

	loadConf := ssg.LoadConfig{
		ContentDir: "content",
		Rules: []ssg.Rule{
			{
				Pattern:   "**/*.html",
				Loader:    metajson.Loader,
				Template:  "baseof.html",
				Transform: ssg.CleanURLs(".html", ".html"),
			},
		},
	}

	pipeline := []ssg.Renderer{
		tmplMacro(fns),
		htmlclean.Render,
		ssg.Must(ssg.NewPageRender("layout", fns)),
		HTMLPretty,
		ssg.WriteOutput("public"),
	}

	pages := []ssg.ContentSourceConfig{}
	if err := ssg.LoadContent(loadConf, &pages); err != nil {
		log.Fatalf("load content: %s", err)
	}

	// Build taxonomy: group content pages by tag.
	byTag := ssg.GroupByStrings(pages, "Tags")

	// Sort tag names for deterministic output.
	tagNames := make([]string, 0, len(byTag))
	for tag := range byTag {
		tagNames = append(tagNames, tag)
	}
	sort.Strings(tagNames)

	// One listing page per tag.
	for _, tag := range tagNames {
		pages = append(pages, ssg.NewPage(
			"tags/"+slug(tag)+"/index.html",
			"tag-list/index.html",
			map[string]any{
				"Title": "Tag: " + tag,
				"Tag":   tag,
				"Pages": byTag[tag],
			},
		))
	}

	// Tag index: precompute name+count structs so the template stays simple.
	tagList := make([]map[string]any, 0, len(tagNames))
	for _, tag := range tagNames {
		tagList = append(tagList, map[string]any{
			"Name":  tag,
			"Count": len(byTag[tag]),
		})
	}
	pages = append(pages, ssg.NewPage(
		"tags/index.html",
		"tag-index/index.html",
		map[string]any{
			"Title": "All Tags",
			"Tags":  tagList,
		},
	))

	if err := ssg.Render(pipeline, pages, nil); err != nil {
		log.Fatalf("render: %s", err)
	}
}
