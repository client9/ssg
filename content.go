package ssg

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// FileWalker returns a Plugin that walks contentDir, matches each file against
// rules in order, and appends the resulting Artifacts to the slice.
//
// For each matched file the Loader is called once. Each Output in the matched
// Rule produces one Artifact, sharing the same parsed metadata but carrying its
// own Pipeline — enabling one-to-many outputs from a single source file.
//
// Files matching no rule, or whose rule has a nil Loader, are skipped.
// Directories prefixed with "." are skipped entirely.
func FileWalker(contentDir string, rules []Rule) Plugin {

	return func(ctx *Context, artifacts *[]Artifact) error {
		return filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("FileWalker: walking %q: %v", path, err)
			}
			if d.IsDir() {
				name := d.Name()
				if name != "." && strings.HasPrefix(d.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}

			relPath := path
			if contentDir != "." {
				relPath = filepath.ToSlash(path[len(contentDir)+1:])
			}

			rule, ok, err := matchRules(rules, relPath)
			if err != nil {
				return fmt.Errorf("FileWalker: bad pattern matching %q: %w", relPath, err)
			}
			if !ok || rule.Loader == nil {
				return nil
			}

			raw, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("FileWalker: reading %s: %w", path, err)
			}

			rawMeta, body, err := rule.Loader(raw)
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			if rawMeta == nil {
				return nil // loader signalled skip
			}

			base := ContentSourceConfig(rawMeta)

			base["Content"] = body
			base["InputFile"] = path
			base["SourcePath"] = relPath
			a := Artifact{
				Meta:     ContentSourceConfig(maps.Clone(map[string]any(base))),
				Pipeline: rule.Pipeline,
			}
			*artifacts = append(*artifacts, a)
			return nil
		})
	}
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
