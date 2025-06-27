package ssg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type SiteConfig struct {
	TemplateDir  string
	BaseTemplate string

	ContentDir string

	Split      ContentSplitter
	Metaparser ParseMeta

	InputExt    string // ".md"
	OutputExt   string // ".html"
	IndexSource string // "index.md"
	IndexDest   string // "index.html"

	Pipeline []Renderer
}

func Main2(config SiteConfig, pages *[]ContentSourceConfig) error {

	// both of these
	//   maybe should be pulled out of function

	ConfigDefault(&config)

	// load in content
	if err := LoadContent(config, pages); err != nil {
		return fmt.Errorf("load content failed: %w", err)
	}

	// TBD: do global site stuff

	for _, p := range *pages {
		// give every page the global config
		p["Site"] = config

		// initial source is in []byte
		source := p["Content"].([]byte)

		if err := MultiRender(config.Pipeline, source, p); err != nil {
			return err
		}
	}

	return nil
}

func LoadContent(config SiteConfig, out *[]ContentSourceConfig) error {

	contentDir := config.ContentDir
	//log.Printf("In content dir: %s", contentDir)

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		//log.Printf("LoadContent: got %s", path)
		// not sure how this works
		if err != nil {
			return fmt.Errorf("LoadContent walkdir error: %v", err)
		}

		// do not look at linux/mac dot dirs
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		//
		if !strings.HasSuffix(path, config.InputExt) {
			return nil
		}

		//log.Printf("LoadContent: reading %s", path)
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("LoadContent: reading page file %s failed: %w", path, err)
		}

		head, body := config.Split(raw)
		// TODO, if head is nil, then we should just copy
		// i.e.  if doc is "{\n}body" that's fine, keep going
		// i.e.  but if doc is "body", then just copy.

		page, err := config.Metaparser(head)
		if err != nil {
			return fmt.Errorf("unable to parse front matter: %v", err)
		}

		if _, ok := page["TemplateName"]; !ok {
			page["TemplateName"] = config.BaseTemplate
		}

		// This name change should be a function

		// have: content/foo/bar/page.sh
		// want: foo/bar/page/index.html
		if _, ok := page["OutputFile"]; !ok {
			s := path[len(contentDir)+1:]
			s = strings.TrimSuffix(s, filepath.Ext(s))

			// in same dir --> index.md --> index.html
			// or make dir --> foo.md --> foo/index.html

			if d.Name() == config.IndexSource && config.IndexDest != "" {
				s += config.OutputExt
			} else {
				s = filepath.Join(s, config.IndexDest)
			}
			//log.Printf("Setting outfile to %s", s)
			page["OutputFile"] = s
		}

		page["Content"] = body
		*out = append(*out, page)
		return nil
	})

	return err
}
