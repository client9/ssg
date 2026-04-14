// Package minify provides a DynStage that minifies HTML, CSS, JS, SVG, and JSON.
//
// The MIME type is derived from the page's OutputFile extension, so no
// per-stage configuration is required. Unknown extensions pass through
// unchanged.
package minify

import (
	"bytes"
	"path/filepath"

	"github.com/client9/ssg"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

// extMime maps file extensions to the MIME types registered with the minifier.
// Using a hardcoded map avoids platform-specific surprises from the OS MIME
// database (e.g. macOS registering .xyz as a chemistry format).
var extMime = map[string]string{
	".html": "text/html",
	".htm":  "text/html",
	".css":  "text/css",
	".js":   "text/javascript",
	".mjs":  "text/javascript",
	".svg":  "image/svg+xml",
	".json": "application/json",
	".xml":  "application/xml",
}

// New returns a DynStage that minifies content based on the output file's
// extension. Supported extensions: .html, .htm, .css, .js, .mjs, .svg,
// .json, .xml. Unrecognised extensions are passed through without modification.
func New() ssg.Stage {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFunc("application/json", json.Minify)
	m.AddFunc("application/xml", xml.Minify)

	return ssg.Step("minify", func(_ *ssg.Context, cfg ssg.ContentSourceConfig, in []byte) ([]byte, error) {
		mimeType := mimeFromFile(cfg.OutputFile())
		if mimeType == "" {
			return in, nil
		}
		var buf bytes.Buffer
		if err := m.Minify(mimeType, &buf, bytes.NewReader(in)); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	})
}

// mimeFromFile returns the MIME type for the given filename based on its
// extension. Returns "" for unrecognised extensions.
func mimeFromFile(name string) string {
	return extMime[filepath.Ext(name)]
}
