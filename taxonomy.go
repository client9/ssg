package ssg

// GroupByString groups artifacts by the string value of a single-value field
// (e.g. "Category"). Artifacts where the field is absent, empty, or not a
// string are skipped.
//
//	byCategory := ssg.GroupByString(artifacts, "Category")
//	for cat, catArtifacts := range byCategory { ... }
func GroupByString(artifacts []Artifact, field string) map[string][]Artifact {
	out := make(map[string][]Artifact)
	for _, a := range artifacts {
		v, ok := a.Meta[field].(string)
		if !ok || v == "" {
			continue
		}
		out[v] = append(out[v], a)
	}
	return out
}

// GroupByStrings groups artifacts by a multi-value string field (e.g. "Tags").
// An artifact appears in the result once for each value it contains.
//
// The field value may be []string, []any (as produced by YAML/JSON parsers),
// or a bare string (a single value written without a list). Empty strings are
// skipped.
//
//	byTag := ssg.GroupByStrings(artifacts, "Tags")
//	for tag, tagArtifacts := range byTag { ... }
func GroupByStrings(artifacts []Artifact, field string) map[string][]Artifact {
	out := make(map[string][]Artifact)
	for _, a := range artifacts {
		for _, v := range toStringSlice(a.Meta[field]) {
			if v == "" {
				continue
			}
			out[v] = append(out[v], a)
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
