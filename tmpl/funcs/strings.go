package funcs

import (
	"strings"
	"text/template"
)

func stringFuncMap() template.FuncMap {
	return template.FuncMap{
		// Case
		"lower": strings.ToLower,
		"upper": strings.ToUpper,

		// Whitespace
		"trim":       strings.TrimSpace,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"trimLeft":   strings.TrimLeft,
		"trimRight":  strings.TrimRight,

		// Search
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"count":     strings.Count,

		// Transform
		"replaceAll": strings.ReplaceAll,
		"repeat":     strings.Repeat,

		// Split / join
		"split":  strings.Split,
		"join":   strings.Join,
		"fields": strings.Fields,

		// Truncate with ellipsis
		"truncate": truncate,
	}
}

// truncate shortens s to at most n runes. If s is longer it is cut and an
// ellipsis ("…") is appended. n includes the ellipsis, so the result is
// always at most n runes long.
//
//	truncate("hello world", 8) → "hello w…"
//	truncate("hi", 8)          → "hi"
func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n == 1 {
		return "…"
	}
	return string(runes[:n-1]) + "…"
}
