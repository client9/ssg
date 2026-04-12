package ssg

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteSitemap generates a sitemap.xml in outDir listing all pages.
// baseURL is prepended to each page's OutputFile to form the full URL
// (trailing slashes on baseURL are trimmed automatically).
// Pages with an empty OutputFile are skipped.
//
// The output follows the sitemap.org protocol so search engines can
// discover all pages on the site.
//
// Example:
//
//	pages := []ssg.ContentSourceConfig{}
//	ssg.LoadContent(conf, &pages)
//	ssg.Render(conf, pages, nil)
//	ssg.WriteSitemap("public", "https://example.com", pages)
func WriteSitemap(outDir, baseURL string, pages []ContentSourceConfig) error {
	baseURL = strings.TrimRight(baseURL, "/")

	type url struct {
		Loc string `xml:"loc"`
	}
	type urlset struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []url    `xml:"url"`
	}

	us := urlset{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, p := range pages {
		out := p.OutputFile()
		if out == "" {
			continue
		}
		// Use forward slashes in URLs regardless of OS path separator.
		loc := baseURL + "/" + filepath.ToSlash(out)
		us.URLs = append(us.URLs, url{Loc: loc})
	}

	dest := filepath.Join(outDir, "sitemap.xml")
	if err := os.MkdirAll(outDir, 0750); err != nil {
		return fmt.Errorf("WriteSitemap: %w", err)
	}
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("WriteSitemap: %w", err)
	}
	defer func() {
		if cerr := f.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	if _, err = fmt.Fprint(f, xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	if err = enc.Encode(us); err != nil {
		return err
	}
	return enc.Close()
}
