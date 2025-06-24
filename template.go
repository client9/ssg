package ssg

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

func NewPageRender(tdir string, fns template.FuncMap) (Renderer, error) {

	// load in all the golang templates
	tmpl, err := templateMap(tdir, fns)
	if err != nil {
		return nil, err
	}

	return func(wr io.Writer, source io.Reader, data any) error {
		s := data.(ContentSourceConfig)

		src, err := io.ReadAll(source)
		if err != nil {
			return err
		}

		// needs to be string for golang text/template
		s["Content"] = string(src)

		return tmpl.ExecuteTemplate(wr, s.TemplateName(), s)
	}, nil
}

func templateMap(root string, fmap template.FuncMap) (TemplateRouter, error) {
	// init
	out := make(TemplateRouter)

	if fmap == nil {
		fmap = template.FuncMap{}
	}

	//
	dirs, err := getDirectories(root)
	if err != nil {
		return nil, err
	}
	for _, d := range dirs {
		t := template.New(d).Funcs(fmap)

		parts := strings.Split(d, string(filepath.Separator))
		for i := 0; i <= len(parts); i++ {
			current := filepath.Join(parts[:i]...)
			templateGlob := filepath.Join(root, current, "*.html")
			if _, err := t.ParseGlob(templateGlob); err != nil {
				return out, err
				// typically empty directory
				//log.Printf("GOT ZERO TEMPLATES")
			}
		}
		out[d] = t
	}
	return out, nil
}

// in the layout directory
//
//	get a list of all paths
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

type TemplateRouter map[string]*template.Template

func (t TemplateRouter) ExecuteTemplate(wr io.Writer, name string, data any) error {
	dir, file := filepath.Split(name)
	dir = strings.TrimSuffix(dir, string(filepath.Separator))
	// fix some asymmetries
	if dir == "" {
		dir = "."
	}
	base, ok := t[dir]
	if !ok {
		return fmt.Errorf("could not file with dir=%q file=%q", dir, file)
	}
	return base.ExecuteTemplate(wr, file, data)
}
