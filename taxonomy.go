package ssg

// GroupByString groups pages by the string value of a single-value field
// (e.g. "Category"). Pages where the field is absent, empty, or not a string
// are skipped.
//
//	byCategory := ssg.GroupByString(pages, "Category")
//	for cat, catPages := range byCategory { ... }
func GroupByString(pages []ContentSourceConfig, field string) map[string][]ContentSourceConfig {
	out := make(map[string][]ContentSourceConfig)
	for _, p := range pages {
		v, ok := p[field].(string)
		if !ok || v == "" {
			continue
		}
		out[v] = append(out[v], p)
	}
	return out
}

// GroupByStrings groups pages by a multi-value string field (e.g. "Tags").
// A page appears in the result once for each value it contains.
//
// The field value may be []string, []any (as produced by YAML/JSON parsers),
// or a bare string (a single value written without a list). Empty strings are
// skipped.
//
//	byTag := ssg.GroupByStrings(pages, "Tags")
//	for tag, tagPages := range byTag { ... }
func GroupByStrings(pages []ContentSourceConfig, field string) map[string][]ContentSourceConfig {
	out := make(map[string][]ContentSourceConfig)
	for _, p := range pages {
		for _, v := range toStringSlice(p[field]) {
			if v == "" {
				continue
			}
			out[v] = append(out[v], p)
		}
	}
	return out
}

// toStringSlice normalises a field value to []string, handling the types
// that frontmatter parsers commonly produce.
func toStringSlice(v any) []string {
	switch val := v.(type) {
	case string:
		return []string{val}
	case []string:
		return val
	case []any:
		out := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}
