package ssg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// LoadConfig holds parameters for LoadContent. It is concerned only with
// finding and parsing content files — it has no role in rendering.
type LoadConfig struct {
	ContentDir   string
	BaseTemplate string

	MetaSplit ContentSplitter
	MetaParser MetaParser

	InputExt    string // e.g. ".md"
	OutputExt   string // e.g. ".html"
	IndexSource string // e.g. "index.md"
	IndexDest   string // e.g. "index.html"
}

// Render renders all pages through pipeline.
//
// globals contains site-wide data (navigation menus, tag indexes, etc.)
// computed after LoadContent. Each entry is merged into the page's
// ContentSourceConfig before rendering. Page frontmatter wins on key
// collision — globals act as defaults, not overrides. globals may be nil.
func Render(pipeline []Renderer, pages []ContentSourceConfig, globals map[string]any) error {
	for _, p := range pages {
		for k, v := range globals {
			if _, exists := p[k]; !exists {
				p[k] = v
			}
		}
		source := p["Content"].([]byte)
		if err := MultiRender(pipeline, source, p); err != nil {
			return fmt.Errorf("%s: %w", p.InputFile(), err)
		}
	}
	return nil
}

// LoadContent walks conf.ContentDir, parses each matching file's frontmatter,
// and appends a ContentSourceConfig to out for each page found.
func LoadContent(conf LoadConfig, out *[]ContentSourceConfig) error {
	contentDir := conf.ContentDir
	if contentDir == "" {
		return fmt.Errorf("ContentDir in config is empty")
	}

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("LoadContent walkdir error @ %q: %v", path, err)
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, conf.InputExt) {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("LoadContent: reading page file %s failed: %w", path, err)
		}

		head, body := conf.MetaSplit(raw)

		page, err := conf.MetaParser(head)
		if err != nil {
			return fmt.Errorf("unable to parse front matter: %v", err)
		}

		if _, ok := page["TemplateName"]; !ok {
			page["TemplateName"] = conf.BaseTemplate
		}

		if _, ok := page["OutputFile"]; !ok {
			s := path[len(contentDir)+1:]
			s = strings.TrimSuffix(s, filepath.Ext(s))

			if d.Name() == conf.IndexSource && conf.IndexDest != "" {
				s += conf.OutputExt
			} else {
				s = filepath.Join(s, conf.IndexDest)
			}
			page["OutputFile"] = s
		}

		page["InputFile"] = path
		page["Content"] = body
		*out = append(*out, page)
		return nil
	})

	return err
}
