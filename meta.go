package ssg

import (
	"bytes"
)

// Splits Input into metadata and the main content/body
type ContentSplitter func(s []byte) ([]byte, []byte)

// ContentSource only needs to know what template to use and
// where the output is going.
type ContentSource interface {
	TemplateName() string
	InputFile() string
	OutputFile() string
}

// ParseMeta parses the front matter and returns a content source
type ParseMeta func(s []byte) (ContentSourceConfig, error)

// Splits a document into a head and body based on various markers.
//

type HeadType struct {
	Name        string
	Prefix      []byte
	Suffix      []byte
	KeepMarkers bool
}

var HeadYaml = HeadType{
	Name:        "yaml",
	Prefix:      []byte("---\n"),
	Suffix:      []byte("\n---\n"),
	KeepMarkers: false,
}
var HeadJson = HeadType{
	Name:        "json",
	Prefix:      []byte("{\n"),
	Suffix:      []byte("\n}\n"),
	KeepMarkers: true,
}
var HeadToml = HeadType{
	Name:        "toml",
	Prefix:      []byte("+++\n"),
	Suffix:      []byte("\n+++\n"),
	KeepMarkers: false,
}
var HeadEmail = HeadType{
	Name:        "email",
	Prefix:      []byte(""),
	Suffix:      []byte("\n\n\n"),
	KeepMarkers: false,
}

func Splitter(head HeadType, s []byte) ([]byte, []byte) {
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

func Joiner(head HeadType, meta []byte, body []byte) []byte {
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
