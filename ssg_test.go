package ssg_test

import (
	"encoding/json"
	"strings"
	"testing"
	"text/template"

	"github.com/npg70/ssg"
)

type siteConfig struct {
	outputDir string
}

func (sconfig siteConfig) OutputDir() string {
	return sconfig.outputDir
}

func TestEmpty(t *testing.T) {

	config := siteConfig{
		outputDir: "",
	}

	tmpl, err := template.New("test").Parse("This is {{.Page.Content}}\n")
	if err != nil {
		t.Errorf("Template init failed: %v", err)
	}

	page := make(ssg.ContentSourceConfig)
	page["TemplateName"] = "test"
	page["OutputFile"] = "/tmp/junk"
	page["Content"] = "content!"

	pages := []ssg.ContentSource{
		page,
	}

	if err := ssg.Execute(config, tmpl, pages); err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	//todo.. read output file.  Check it matches.
}

func TestSimpleYamlContent(t *testing.T) {

	config := siteConfig{
		outputDir: "",
	}

	tmpl, err := template.New("test").Parse("This is {{.Page.Content}}\n")
	if err != nil {
		t.Errorf("Template init failed: %v", err)
	}

	doc := strings.TrimSpace(`
{
	"TemplateName": "test",
	"OutputFile": "/tmp/junk2"
}
Multi
  Line
    Content
`)
	cs := ssg.ContentSplitter{}
	cs.Register(ssg.HeadJson)
	htype, head, body := cs.Split([]byte(doc))
	if htype != "json" {
		t.Errorf("Expected JSON sample: got %q", htype)
	}

	page := make(ssg.ContentSourceConfig)
	if err := json.Unmarshal(head, &page); err != nil {
		t.Errorf("Unable to un-yaml: %v", err)
	}
	page["Content"] = body

	pages := []ssg.ContentSource{
		page,
	}

	if err := ssg.Execute(config, tmpl, pages); err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	//todo.. read output file.  Check it matches.
}
