package ssg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// LoadConfig holds parameters for LoadContent. It is concerned only with
// finding and parsing content files — it has no role in rendering.
type LoadConfig struct {
	ContentDir string

	// Rules are tried in order against each file's relative path.
	// The first matching Rule's Loader is called. Files that match no
	// rule are skipped. Use doublestar glob syntax: "**/*.md", "posts/*.html".
	Rules []Rule
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
		source, _ := p["Content"].([]byte)
		if err := MultiRender(pipeline, source, p); err != nil {
			return fmt.Errorf("%s: %w", p.InputFile(), err)
		}
	}
	return nil
}

// LoadContent walks conf.ContentDir, matches each file against conf.Rules in
// order, and appends the resulting pages to out. Files matching no rule are
// skipped. Directories prefixed with "." are skipped entirely.
func LoadContent(conf LoadConfig, out *[]ContentSourceConfig) error {
	if conf.ContentDir == "" {
		return fmt.Errorf("ContentDir in config is empty")
	}

	return filepath.WalkDir(conf.ContentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("LoadContent: walking %q: %v", path, err)
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		relPath := filepath.ToSlash(path[len(conf.ContentDir)+1:])

		rule, ok, err := matchRules(conf.Rules, relPath)
		if err != nil {
			return fmt.Errorf("LoadContent: bad pattern matching %q: %w", relPath, err)
		}
		if !ok || rule.Loader == nil {
			return nil // no rule matched, or rule explicitly skips
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("LoadContent: reading %s: %w", path, err)
		}

		rawMeta, body, err := rule.Loader(raw)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if rawMeta == nil {
			return nil // loader signalled skip
		}

		page := ContentSourceConfig(rawMeta)
		page["Content"] = body

		// Fill OutputFile from Rule.Transform if frontmatter didn't set it.
		if page.OutputFile() == "" {
			var outputFile string
			if rule.Transform != nil {
				outputFile = rule.Transform(relPath)
				if outputFile == "" {
					return nil // transform signalled skip
				}
			} else {
				outputFile = relPath
			}
			page["OutputFile"] = outputFile
		}

		// Fill TemplateName from Rule.Template if frontmatter didn't set it.
		if page.TemplateName() == "" && rule.Template != "" {
			page["TemplateName"] = rule.Template
		}

		page["InputFile"] = path
		*out = append(*out, page)
		return nil
	})
}

// matchRules returns the first Rule whose Pattern matches relPath.
// Returns an error only if a pattern is malformed.
func matchRules(rules []Rule, relPath string) (Rule, bool, error) {
	for _, r := range rules {
		ok, err := doublestar.Match(r.Pattern, relPath)
		if err != nil {
			return Rule{}, false, fmt.Errorf("pattern %q: %w", r.Pattern, err)
		}
		if ok {
			return r, true, nil
		}
	}
	return Rule{}, false, nil
}
