package ssg

// Rule pairs a doublestar glob pattern with a MetaLoader and a Pipeline.
// FileWalker tries rules in order; the first pattern that matches the file's
// relative path wins. Files that match no rule are skipped.
// A nil Loader skips the file without reading it.
//
//	Rule{
//	    Pattern:  "**/*.md",
//	    Loader:   metayaml.Loader,
//	    Pipeline: ssg.NewPipeline("post",
//	        ssg.SetOutputFile(ssg.CleanURLs(".md", ".html")),
//	        markdown.New(),
//	        ssg.Must(ssg.NewPageRender("layout", fns)),
//	        ssg.WriteOutput,
//	    ),
//	}
//	Rule{Pattern: "**/_*"} // nil Loader: skip draft files
type Rule struct {
	Pattern  string
	Loader   MetaLoader
	Pipeline Pipeline
}

// Passthrough is a MetaLoader that returns the raw file bytes as body with
// empty metadata. Use it for assets (images, CSS, JS) that should be copied
// to the output directory unchanged.
var Passthrough MetaLoader = func(raw []byte) (map[string]any, []byte, error) {
	return map[string]any{}, raw, nil
}

// Skip is a MetaLoader that unconditionally skips the file.
// A nil Loader in a Rule has the same effect; Skip makes the intent explicit.
//
//	Rule{Pattern: "**/_*", Loader: ssg.Skip}
var Skip MetaLoader = func(_ []byte) (map[string]any, []byte, error) {
	return nil, nil, nil
}
