package ssg

import (
	"path/filepath"
	"strings"
)

// PathTransformer maps a relative input path to a relative output path.
// The input path is relative to ContentDir.
// Returning an empty string skips the file.
//
// Use CleanURLs or UglyURLs for the common cases, and wrap with modifiers
// like SlugNormalize to compose behaviour:
//
//	ssg.SlugNormalize(ssg.CleanURLs(".md", ".html"))
type PathTransformer func(relPath string) string

// CleanURLs returns a PathTransformer that produces directory-style output
// paths, making pages accessible without a file extension in the browser.
//
//	foo.md       → foo/index.html     (/foo or /foo/)
//	bar/baz.md   → bar/baz/index.html
//	index.md     → index.html         (root index, not index/index.html)
func CleanURLs(inputExt, outputExt string) PathTransformer {
	return func(relPath string) string {
		if filepath.Ext(relPath) != inputExt {
			return ""
		}
		base := strings.TrimSuffix(relPath, inputExt)
		// "index" at any level stays as index.html, not index/index.html
		if filepath.Base(base) == "index" {
			return filepath.ToSlash(base + outputExt)
		}
		return filepath.ToSlash(filepath.Join(base, "index"+outputExt))
	}
}

// UglyURLs returns a PathTransformer that maps input files directly to
// same-named output files, replacing the extension.
//
//	foo.md       → foo.html
//	bar/baz.md   → bar/baz.html
func UglyURLs(inputExt, outputExt string) PathTransformer {
	return func(relPath string) string {
		if filepath.Ext(relPath) != inputExt {
			return ""
		}
		return filepath.ToSlash(strings.TrimSuffix(relPath, inputExt) + outputExt)
	}
}

// SlugNormalize returns a PathTransformer that lowercases the filename and
// replaces spaces and underscores with hyphens before passing to next.
//
//	"Foo Bar.md"  → "foo-bar.md"  → (next applies)
//	"my_post.md"  → "my-post.md"  → (next applies)
func SlugNormalize(next PathTransformer) PathTransformer {
	return func(relPath string) string {
		dir := filepath.Dir(relPath)
		base := filepath.Base(relPath)
		base = strings.ToLower(base)
		base = strings.ReplaceAll(base, " ", "-")
		base = strings.ReplaceAll(base, "_", "-")
		normalized := filepath.Join(dir, base)
		return next(normalized)
	}
}

// SetTemplateName returns a DynStage that sets TemplateName in cfg to name,
// but only if frontmatter hasn't already set it. Content passes through
// unchanged, making this a metadata-only pipeline step.
func SetTemplateName(name string) Stage {
	return Step("set-template-name", func(_ *Context, cfg ContentSourceConfig, in any) (any, error) {
		if cfg.TemplateName() == "" {
			cfg["TemplateName"] = name
		}
		return in, nil
	})
}

// SetOutputFile returns a DynStage that applies transform to the artifact's
// SourcePath and stores the result as OutputFile in cfg. Content passes
// through unchanged, making this a metadata-only pipeline step.
// If transform returns "" the OutputFile is left unchanged.
func SetOutputFile(transform PathTransformer) Stage {
	return Step("set-output-file", func(_ *Context, cfg ContentSourceConfig, in any) (any, error) {
		relPath := cfg.SourcePath()
		if relPath == "" {
			relPath = cfg.InputFile()
		}
		if out := transform(relPath); out != "" {
			cfg["OutputFile"] = out
		}
		return in, nil
	})
}
