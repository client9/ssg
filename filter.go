package ssg

// FilterArtifacts returns a Plugin that retains only the artifacts for which
// fn returns true. Use it to remove draft pages, future-dated posts, etc.
//
// Example — exclude draft pages:
//
//	ssg.FilterArtifacts(func(meta ssg.ContentSourceConfig) bool {
//	    draft, _ := meta["draft"].(bool)
//	    return !draft
//	})
func FilterArtifacts(fn func(ContentSourceConfig) bool) Plugin {
	return func(ctx *Context, artifacts *[]Artifact) error {
		out := (*artifacts)[:0]
		for _, a := range *artifacts {
			if fn(a.Meta) {
				out = append(out, a)
			}
		}
		*artifacts = out
		return nil
	}
}
