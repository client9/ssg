package ssg

import "bytes"

// Splits a document into a head and body based on various markers.
// Other implimentations parse the head and/or body.
// This just splits the two appropriately and does not parse either.
// BYOP - Bring Your Own Parser.

// if you never use TOML, why spend any cycles checking for it, and why
// have it as a dependency?
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

type ContentSplitter struct {
	formats []HeadType
}

func (cs *ContentSplitter) Register(m HeadType) {
	cs.formats = append(cs.formats, m)
}

func (cs *ContentSplitter) Split(s []byte) (string, []byte, []byte) {
	for _, head := range cs.formats {
		if bytes.HasPrefix(s, head.Prefix) {
			plen := len(head.Prefix)
			if idx := bytes.Index(s[plen:], head.Suffix); idx != -1 {
				if head.KeepMarkers {
					pt := plen + idx + len(head.Suffix)
					return head.Name, s[:pt], s[pt:]
				}
				return head.Name, s[plen : plen+idx], s[plen+idx+len(head.Suffix):]
			}
		}
	}
	return "", nil, s
}
