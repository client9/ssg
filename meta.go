package ssg

import (
	"strings"
)

// Splits a document into a head and body based on various markers.
// Other implimentations parse the head and/or body.
// This just splits the two appropriately and does not parse either.
// BYOP - Bring Your Own Parser.

// if you never use TOML, why spend any cycles checking for it, and why
// has a dependency?
//

type HeadType struct {
	Name        string
	Prefix      string
	Suffix      string
	KeepMarkers bool
}

var HeadYaml = HeadType{
	Name:        "yaml",
	Prefix:      "---\n",
	Suffix:      "\n---\n",
	KeepMarkers: false,
}
var HeadJson = HeadType{
	Name:        "json",
	Prefix:      "{\n",
	Suffix:      "\n}\n",
	KeepMarkers: true,
}
var HeadToml = HeadType{
	Name:        "toml",
	Prefix:      "+++\n",
	Suffix:      "\n+++\n",
	KeepMarkers: false,
}
var HeadEmail = HeadType{
	Name:        "email",
	Prefix:      "",
	Suffix:      "\n\n\n",
	KeepMarkers: false,
}

type ContentSplitter struct {
	formats []HeadType
}

func (cs *ContentSplitter) Register(m HeadType) {
	cs.formats = append(cs.formats, m)
}

func (cs *ContentSplitter) Split(s string) (string, string, string) {
	for _, head := range cs.formats {
		if strings.HasPrefix(s, head.Prefix) {
			plen := len(head.Prefix)
			if idx := strings.Index(s[plen:], head.Suffix); idx != -1 {
				if head.KeepMarkers {
					pt := plen + idx + len(head.Suffix)
					return head.Name, s[:pt], s[pt:]
				}
				return head.Name, s[plen : plen+idx], s[plen+idx+len(head.Suffix):]
			}
		}
	}
	return "", "", s
}
