package ssg

// FileLoader interprets a single file into a ContentSourceConfig.
// Returning (nil, nil) signals that the file should be skipped.
type FileLoader func(relPath string, raw []byte) (ContentSourceConfig, error)

// Rule pairs a doublestar glob pattern with a FileLoader.
// LoadContent tries rules in order; the first pattern that matches the file's
// relative path wins. Files that match no rule are skipped.
//
//	Rule{"**/*.md",  FrontmatterLoader(...)}
//	Rule{"**/*.css", PassthroughLoader()}
type Rule struct {
	Pattern string
	Loader  FileLoader
}

// FrontmatterLoader returns a FileLoader that uses meta to parse each file's
// frontmatter and body, then resolves the output path and template name.
// If transform returns "" for a given path the file is skipped.
func FrontmatterLoader(meta MetaLoader, baseTemplate string, transform PathTransformer) FileLoader {
	return func(relPath string, raw []byte) (ContentSourceConfig, error) {
		page, body, err := meta(raw)
		if err != nil {
			return nil, err
		}

		outputFile := page.OutputFile()
		if outputFile == "" {
			if transform != nil {
				outputFile = transform(relPath)
			}
			if outputFile == "" {
				return nil, nil // skip
			}
		}

		templateName := page.TemplateName()
		if templateName == "" {
			templateName = baseTemplate
		}

		p := NewPage(outputFile, templateName, page)
		p["Content"] = body
		return p, nil
	}
}

// PassthroughLoader returns a FileLoader that copies each file to the output
// directory unchanged, preserving its relative path. Use this for assets
// (images, CSS, JS) that live alongside content files.
func PassthroughLoader() FileLoader {
	return func(relPath string, raw []byte) (ContentSourceConfig, error) {
		p := NewPage(relPath, "", nil)
		p["Content"] = raw
		return p, nil
	}
}
