package ssg

import "log"

// Context holds site-wide state shared across all Plugins and Renderers.
// It is readable and writable by plugins — for example, a nav-building plugin
// can compute navigation and store it in Globals for templates to consume.
type Context struct {
	Globals   map[string]any
	OutputDir string
	Logger    *log.Logger
}
