package ssg

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

// NewPageRender loads HTML templates from tdir and returns a DynStage that
// wraps each page's rendered content in the appropriate layout template.
//
// Template selection: each page's ContentSourceConfig must contain a
// "TemplateName" key (e.g. "blog/single.html"). The directory portion routes
// to the right template set; the filename selects the template within it.
//
// Content injection: the stage reads the pipeline's current output, stores
// it as the string value of page["Content"], then executes the named template
// with the full ContentSourceConfig as its data. Templates access the body via
// {{.Content}} and other page metadata via {{.Title}}, {{.Date}}, etc.
//
// Template inheritance: templates in a parent directory are available to all
// templates in child directories. A template in layout/blog/ can call
// {{template "base.html" .}} because base.html is parsed into the same set.
// See TemplateRouter for details.
//
// Block override constraint: all templates in the same directory share one
// template set, so only one template per directory may use {{define "main"}}
// (or any other block name). If two sibling templates both define the same
// block, Go's template engine will error. The solution is to give each
// template that overrides a block its own subdirectory:
//
//	layout/
//	  baseof.html       ← defines {{block "main" .}}
//	  post/
//	    index.html      ← {{define "main"}} for post pages
//	  tag-list/
//	    index.html      ← {{define "main"}} for tag listing pages
//	  tag-index/
//	    index.html      ← {{define "main"}} for the tag index page
//
// Each subdirectory gets an isolated template set that inherits baseof.html
// from the parent but does not share its set with siblings.
//
// fns is an optional map of additional template functions made available to
// all templates. Pass nil for no extra functions.
func NewPageRender(tdir string, fns template.FuncMap) (Stage, error) {
	tmpl, err := templateMap(tdir, fns)
	if err != nil {
		return nil, err
	}

	return Step("page-render", func(_ *Context, cfg ContentSourceConfig, in []byte) ([]byte, error) {
		// Store the rendered body as template.HTML so html/template does not
		// escape it — the content is already trusted, rendered markup.
		// All other page metadata (Title, Author, etc.) is auto-escaped.
		cfg["Content"] = template.HTML(in) //nolint:gosec

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, cfg.TemplateName(), cfg); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}), nil
}

// templateMap walks the layout directory rooted at root and builds a
// TemplateRouter.
//
// For each subdirectory (including "."), it creates a *template.Template that
// has parsed all *.html files from the root directory down to and including
// that subdirectory. This means a template at any level can call templates
// defined in ancestor directories.
//
// Example layout tree:
//
//	layout/
//	  base.html          ← available to all template sets
//	  blog/
//	    single.html      ← can call {{template "base.html" .}}
//	    list/
//	      index.html     ← can call base.html and any blog/*.html template
//
// The resulting TemplateRouter has keys ".", "blog", and "blog/list".
func templateMap(root string, fmap template.FuncMap) (TemplateRouter, error) {
	out := make(TemplateRouter)

	if fmap == nil {
		fmap = template.FuncMap{}
	}

	dirs, err := getDirectories(root)
	if err != nil {
		return nil, err
	}

	for _, d := range dirs {
		t := template.New(d).Funcs(fmap)

		// Load templates from root down to this directory so that ancestor
		// templates are available within child template sets.
		parts := strings.Split(d, string(filepath.Separator))
		for i := 0; i <= len(parts); i++ {
			current := filepath.Join(parts[:i]...)
			templateGlob := filepath.Join(root, current, "*.html")
			if _, err := t.ParseGlob(templateGlob); err != nil {
				return out, err
			}
		}
		out[d] = t
	}
	return out, nil
}

// getDirectories returns all directory paths under root as paths relative to
// root, including root itself as ".".
func getDirectories(root string) ([]string, error) {
	out := []string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir error at %q: %v", path, err)
		}
		if d.IsDir() {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				panic("WTF")
			}
			out = append(out, rel)
		}
		return nil
	})
	return out, err
}

// TemplateRouter maps layout directory paths (relative to the layout root) to
// their compiled template sets. Keys are directory paths such as ".", "blog",
// or "blog/posts".
//
// Each template set contains all templates parsed from the root down to that
// directory, so child sets inherit parent templates.
type TemplateRouter map[string]*template.Template

// ExecuteTemplate routes a template name to the correct set and executes it.
//
// The name is expected to be a slash-separated path where the final component
// is the template filename and everything before it is the directory key.
// For example:
//
//	"base.html"           → directory ".",  template "base.html"
//	"blog/single.html"    → directory "blog", template "single.html"
//	"blog/posts/page.html"→ directory "blog/posts", template "page.html"
func (t TemplateRouter) ExecuteTemplate(wr io.Writer, name string, data any) error {
	dir, file := filepath.Split(name)
	dir = strings.TrimSuffix(dir, string(filepath.Separator))
	if dir == "" {
		dir = "."
	}

	base, ok := t[dir]
	if !ok {
		return fmt.Errorf("no templates loaded for directory %q (from template name %q)", dir, name)
	}
	return base.ExecuteTemplate(wr, file, data)
}
