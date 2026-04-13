package shortcode

import "strings"

// NamedFunc is a handler that receives args as a name→value map instead of a
// positional slice. Use MakeNamed or RegisterNamed to adapt it for use with
// a Context.
type NamedFunc func(ctx *Context, named map[string]string, body string) string

// MakeNamed wraps a NamedFunc as a HandlerFunc. paramNames assigns names to
// positional arguments in order; named arguments (key=value) are passed
// through directly. Extra positional args beyond len(paramNames) are dropped.
//
//	fn := MakeNamed(myFunc, "src", "alt")
//	// $cmd[photo.jpg "a photo"]  →  {"src":"photo.jpg", "alt":"a photo"}
//	// $cmd[src=photo.jpg alt="a photo"]  →  same
func MakeNamed(fn NamedFunc, paramNames ...string) HandlerFunc {
	return func(ctx *Context, args []string, body string) string {
		named := make(map[string]string, len(args))
		posIdx := 0
		for _, arg := range args {
			k, v, ok := strings.Cut(arg, "=")
			if ok {
				named[k] = v
			} else if posIdx < len(paramNames) {
				named[paramNames[posIdx]] = arg
				posIdx++
			}
		}
		return fn(ctx, named, body)
	}
}

// RegisterNamed registers a NamedFunc under name, wrapping it with MakeNamed.
func (c *Context) RegisterNamed(name string, fn NamedFunc, paramNames ...string) {
	c.Tags[name] = MakeNamed(fn, paramNames...)
}
