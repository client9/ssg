package ssg

// MetaLoader parses raw file bytes into frontmatter metadata and body content.
// Returning a nil map signals that the file should be skipped.
// If the file has no recognisable frontmatter, return an empty map and the
// full raw bytes as body.
//
// The return type is map[string]any rather than ContentSourceConfig so that
// loader implementations have no dependency on the ssg module.
// LoadContent casts the map to ContentSourceConfig before further processing.
type MetaLoader func(raw []byte) (map[string]any, []byte, error)
