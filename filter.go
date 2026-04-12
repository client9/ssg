package ssg

// FilterPages returns a new slice containing only the pages for which fn
// returns true. The original slice is not modified.
//
// Example — exclude draft pages:
//
//	pages = ssg.FilterPages(pages, func(p ssg.ContentSourceConfig) bool {
//	    draft, _ := p["draft"].(bool)
//	    return !draft
//	})
//
// Example — exclude future-dated pages (assuming Date is a time.Time):
//
//	now := time.Now()
//	pages = ssg.FilterPages(pages, func(p ssg.ContentSourceConfig) bool {
//	    date, ok := p["Date"].(time.Time)
//	    return !ok || !date.After(now)
//	})
func FilterPages(pages []ContentSourceConfig, fn func(ContentSourceConfig) bool) []ContentSourceConfig {
	out := make([]ContentSourceConfig, 0, len(pages))
	for _, p := range pages {
		if fn(p) {
			out = append(out, p)
		}
	}
	return out
}
