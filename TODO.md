# TODO

## Missing features

### Dev server
Not needed. For local preview, use:
```
python3 -m http.server --directory public --bind 127.0.0.1 1313
```

### Pagination
Split a `[]ContentSourceConfig` into multiple output pages (`/page/1/`, `/page/2/`).
Builds on top of taxonomy — paginated tag lists, paginated archives.

```go
type Paginated struct {
    Items   []ContentSourceConfig
    Number  int    // 1-based
    Total   int
    HasPrev bool
    HasNext bool
    PrevURL string
    NextURL string
}

func Paginate(pages []ContentSourceConfig, baseURL string, pageSize int) []Paginated
```

The caller creates one synthetic page per `Paginated` struct.

### Sitemap lastmod + RSS / Atom feed
These share a common dependency: extracting a `time.Time` from a page's `Date`
frontmatter field. Improve sitemap first, then RSS is mostly the same work.

**Sitemap improvement** (`sitemap.go`):
- Add `<lastmod>` to each `<url>` entry using the page's `Date` field (omitted if absent).
- Format as `"2006-01-02"` (date-only is valid per sitemap.org spec).

**RSS feed** (new `feed.go`):
```go
func WriteRSSFeed(outDir, baseURL, title, description string, pages []ContentSourceConfig) error
```
- Pages should be pre-sorted by `Date` descending and pre-filtered by the caller.
- Uses the same `Date` field as sitemap lastmod.

Convention: the date frontmatter field is `Date` (matches Hugo, Jekyll, Eleventy).

### Related content
Find pages sharing the most values in a given field (e.g. tags).

```go
func RelatedPages(page ContentSourceConfig, all []ContentSourceConfig, field string, n int) []ContentSourceConfig
```

---

## Future Improvements

### Parallel rendering
- Each page is independent; `Render` could run pages concurrently
- The per-page pipeline (buffer-swap) stays sequential, but pages themselves can fan out

### Asset passthrough
- Parallel copy: fan out file copies with a worker pool; same walk, same API
- Skip unchanged files: compare source/dest `os.Stat` mtime before copying;
  fall back to size comparison if mtimes are unreliable (e.g. after a git checkout)
