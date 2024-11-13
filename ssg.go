package ssg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// not sure if this is needed as an interface
type SiteConfig interface {
	OutputDir() string
}

// PageTemplate --
type PageTemplate interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

// ContentSource only needs to know what template to use and
// where the output is going.
type ContentSource interface {
	TemplateName() string
	InputFile() string
	OutputFile() string
}
type ContentSourceConfig map[string]any

func (csc ContentSourceConfig) TemplateName() string {
	if val, ok := csc["TemplateName"]; ok {
		return val.(string)
	}
	return ""
}
func (csc ContentSourceConfig) OutputFile() string {
	if val, ok := csc["OutputFile"]; ok {
		return val.(string)
	}
	return ""
}
func (csc ContentSourceConfig) InputFile() string {
	if val, ok := csc["InputFile"]; ok {
		return val.(string)
	}
	return ""
}
func (csc ContentSourceConfig) Get(key string) string {
	if val, ok := csc[key]; ok {
		if sval, ok := val.(string); ok {
			return sval
		}
	}
	return ""
}

type TemplateData struct {
	Site SiteConfig
	Page ContentSource
}

// Create your templates
// Create list of page sources
// Compute useful global stuff
// For each page source
// look up template...and execute

// TODO: Add virtual filesystem
// TODO: config
func Execute(sconfig SiteConfig,
	tpl PageTemplate,
	sources []ContentSource) error {

	for _, s := range sources {
		err := executeOne(sconfig, tpl, s)
		if err != nil {
			return fmt.Errorf("input %s failed: %w", s.InputFile(), err)
		}
	}
	return nil
}

func executeOne(sconfig SiteConfig,
	tpl PageTemplate,
	s ContentSource) error {

	outdir := sconfig.OutputDir()
	// make directory
	fullpath := filepath.Join(outdir, s.OutputFile())
	dir := filepath.Dir(fullpath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	// open file
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	data := TemplateData{
		Site: sconfig,
		Page: s,
	}
	// create template input data
	// TODO: swap "s" with a bigger input data struct
	if err := tpl.ExecuteTemplate(f, s.TemplateName(), data); err != nil {
		return err
	}
	return nil
}
