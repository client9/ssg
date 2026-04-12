// Package funcs provides a standard template.FuncMap for use with Go's
// text/template and html/template packages.
//
// All functions follow Go stdlib argument order: the primary string or value
// is the first argument. This matches what Go programmers expect from direct
// calls and avoids the confusion of pipeline-optimized argument order.
//
// Usage:
//
//	import "github.com/client9/ssg/tmpl/funcs"
//
//	t := template.New("page").Funcs(funcs.FuncMap())
//
// To combine with your own functions, use Merge:
//
//	fns := funcs.Merge(funcs.FuncMap(), template.FuncMap{
//	    "myFunc": myFunc,
//	})
//	t := template.New("page").Funcs(fns)
package funcs
