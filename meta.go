package ssg

// MetaLoader parses raw file bytes into frontmatter metadata and body content.
// If the file has no recognisable frontmatter prefix, it should return an empty
// ContentSourceConfig and the full raw bytes as body.
type MetaLoader func(raw []byte) (meta ContentSourceConfig, body []byte, err error)
