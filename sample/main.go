package main

import (
	"fmt"
	"html"
	"log"
	"sort"
	"strings"
	"text/template"

	"github.com/client9/ssg"
	metajsonyaml "github.com/client9/ssg/meta/jsonyaml"
	"github.com/client9/ssg/render/htmlclean"
	"github.com/yosssi/gohtml"
)

// {{ elink "https://.../" "name" }} creates an <a> that opens in a new window.
func elink(href string, body string) string {
	return fmt.Sprintf("<a href=%q target=_blank>%s</a>",
		html.EscapeString(href),
		html.EscapeString(body))
}

var HTMLPretty = ssg.Step("html-pretty", func(_ *ssg.Context, _ ssg.ContentSourceConfig, in []byte) ([]byte, error) {
	return gohtml.FormatBytes(in), nil
})

// slug converts a tag name to a URL-safe string.
func slug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func tmplMacro(funcs template.FuncMap) ssg.Stage {
	t := template.New("_macro")
	if funcs != nil {
		t = t.Funcs(funcs)
	}
	return ssg.Step("tmpl-macro", func(_ *ssg.Context, _ ssg.ContentSourceConfig, in []byte) ([]byte, error) {
		var err error
		t, err = t.Parse(string(in))
		if err != nil {
			return nil, err
		}
		var buf strings.Builder
		if err = t.Execute(&buf, nil); err != nil {
			return nil, err
		}
		return []byte(buf.String()), nil
	})
}

func main() {
	fns := template.FuncMap{
		"upper": strings.ToUpper,
		"elink": elink,
	}

	ctx := &ssg.Context{
		OutputDir: "public",
		Logger:    log.Default(),
	}

	htmlPipeline := ssg.NewPipeline("html",
		ssg.SetOutputFile(ssg.CleanURLs(".html", ".html")),
		ssg.SetTemplateName("baseof.html"),
		tmplMacro(fns),
		htmlclean.Render,
		ssg.Must(ssg.NewPageRender("layout", fns)),
		HTMLPretty,
		ssg.WriteOutput,
	)

	rules := []ssg.Rule{
		{
			Pattern:  "**/*.html",
			Loader:   metajsonyaml.Loader,
			Pipeline: htmlPipeline,
		},
	}

	var artifacts []ssg.Artifact
	plugins := []ssg.Plugin{
		ssg.FileWalker("content", rules),
		buildTagPages(fns),
	}
	for _, p := range plugins {
		if err := p(ctx, &artifacts); err != nil {
			log.Fatalf("plugin: %s", err)
		}
	}

	if err := ssg.Render(ctx, &artifacts); err != nil {
		log.Fatalf("render: %s", err)
	}
}

// buildTagPages returns a Plugin that generates tag listing and index pages.
func buildTagPages(fns template.FuncMap) ssg.Plugin {
	tagPipeline := ssg.NewPipeline("tag",
		ssg.Must(ssg.NewPageRender("layout", fns)),
		HTMLPretty,
		ssg.WriteOutput,
	)

	return func(ctx *ssg.Context, artifacts *[]ssg.Artifact) error {
		byTag := ssg.GroupByStrings(*artifacts, "Tags")

		tagNames := make([]string, 0, len(byTag))
		for tag := range byTag {
			tagNames = append(tagNames, tag)
		}
		sort.Strings(tagNames)

		// One listing page per tag.
		for _, tag := range tagNames {
			*artifacts = append(*artifacts, ssg.NewPage(
				"tags/"+slug(tag)+"/index.html",
				"tag-list/index.html",
				map[string]any{
					"Title": "Tag: " + tag,
					"Tag":   tag,
					"Pages": metaSlice(byTag[tag]),
				},
				tagPipeline,
			))
		}

		// Tag index.
		tagList := make([]map[string]any, 0, len(tagNames))
		for _, tag := range tagNames {
			tagList = append(tagList, map[string]any{
				"Name":  tag,
				"Count": len(byTag[tag]),
			})
		}
		*artifacts = append(*artifacts, ssg.NewPage(
			"tags/index.html",
			"tag-index/index.html",
			map[string]any{
				"Title": "All Tags",
				"Tags":  tagList,
			},
			tagPipeline,
		))

		return nil
	}
}

// metaSlice extracts the Meta map from each artifact for use in templates.
func metaSlice(artifacts []ssg.Artifact) []ssg.ContentSourceConfig {
	out := make([]ssg.ContentSourceConfig, len(artifacts))
	for i, a := range artifacts {
		out[i] = a.Meta
	}
	return out
}
