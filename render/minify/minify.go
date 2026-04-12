// Package minify provides a Renderer that minifies HTML, CSS, JS, SVG, and JSON.
//
// The MIME type is derived from the page's OutputFile extension, so no
// per-renderer configuration is required. Unknown extensions pass through
// unchanged.
package minify

import (
	"io"
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

// New returns a Renderer that minifies content based on the output file's
// extension. Supported extensions: .html, .htm, .css, .js, .mjs, .svg,
// .json, .xml. Unrecognised extensions are passed through without modification.
func New() ssg.Renderer {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFunc("application/json", json.Minify)
	m.AddFunc("application/xml", xml.Minify)

	return func(wr io.Writer, src io.Reader, data any) error {
		cs := data.(ssg.ContentSourceConfig)
		mimeType := mimeFromFile(cs.OutputFile())
		if mimeType == "" {
			_, err := io.Copy(wr, src)
			return err
		}
		return m.Minify(mimeType, wr, src)
	}
}

// mimeFromFile returns the MIME type for the given filename based on its
// extension. Returns "" for unrecognised extensions.
func mimeFromFile(name string) string {
	return extMime[filepath.Ext(name)]
}
