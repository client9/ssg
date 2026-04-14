package ssg

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteSitemap generates a sitemap.xml in outDir listing all artifacts.
// baseURL is prepended to each artifact's OutputFile to form the full URL
// (trailing slashes on baseURL are trimmed automatically).
// Artifacts with an empty OutputFile are skipped.
func WriteSitemap(outDir, baseURL string, artifacts []Artifact) error {
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
	for _, a := range artifacts {
		out := a.Meta.OutputFile()
		if out == "" {
			continue
		}
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
