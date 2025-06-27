package ssg

import (
	"bytes"
	"encoding/json"
)

// Splits Input into metadata and the main content/body
type ContentSplitter func(s []byte) ([]byte, []byte)

// MetaParser parses the front matter and returns a content source
type MetaParser func(s []byte) (ContentSourceConfig, error)

// ParseMetaJson is a default parser, that reads front matter as
// JSON and returns a map[string]any type.
func MetaParseJson(s []byte) (ContentSourceConfig, error) {
	meta := ContentSourceConfig{}
	if err := json.Unmarshal(s, &meta); err != nil {
		return nil, err
	}
	return meta, nil
}

// Splits a document into a head and body based on various markers.
//

type MetaHeadType struct {
	Name        string
	Prefix      []byte
	Suffix      []byte
	KeepMarkers bool
}

var MetaHeadYaml = MetaHeadType{
	Name:        "yaml",
	Prefix:      []byte("---\n"),
	Suffix:      []byte("\n---\n"),
	KeepMarkers: false,
}
var MetaHeadJson = MetaHeadType{
	Name:        "json",
	Prefix:      []byte("{\n"),
	Suffix:      []byte("\n}\n"),
	KeepMarkers: true,
}
var MetaHeadToml = MetaHeadType{
	Name:        "toml",
	Prefix:      []byte("+++\n"),
	Suffix:      []byte("\n+++\n"),
	KeepMarkers: false,
}
var MetaHeadEmail = MetaHeadType{
	Name:        "email",
	Prefix:      []byte(""),
	Suffix:      []byte("\n\n\n"),
	KeepMarkers: false,
}

func MetaSplitEmail(s []byte) ([]byte, []byte) {
	return Splitter(MetaHeadEmail, s)
}

func MetaSplitJson(s []byte) ([]byte, []byte) {
	return Splitter(MetaHeadJson, s)
}

func MetaSplitToml(s []byte) ([]byte, []byte) {
	return Splitter(MetaHeadYaml, s)
}

func MetaSplitYaml(s []byte) ([]byte, []byte) {
	return Splitter(MetaHeadYaml, s)
}

func Splitter(head MetaHeadType, s []byte) ([]byte, []byte) {
	if !bytes.HasPrefix(s, head.Prefix) {
		return nil, s
	}
	plen := len(head.Prefix)
	if idx := bytes.Index(s[plen:], head.Suffix); idx != -1 {
		if head.KeepMarkers {
			pt := plen + idx + len(head.Suffix)
			return s[:pt], s[pt:]
		}
		return s[plen : plen+idx], s[plen+idx+len(head.Suffix):]
	}
	return nil, nil
}

func Joiner(head MetaHeadType, meta []byte, body []byte) []byte {
	out := []byte{}
	if !head.KeepMarkers {
		out = append(out, head.Prefix...)
	}
	out = append(out, bytes.TrimSpace(meta)...)

	if !head.KeepMarkers {
		out = append(out, head.Suffix...)
	} else {
		out = append(out, byte('\n'))
	}

	out = append(out, bytes.TrimSpace(body)...)
	return out
}
