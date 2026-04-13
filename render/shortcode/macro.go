// Package tf2 implements an embedded macro system for plain-text documents.
//
// Macros take the following forms:
//
//	$cmd
//	$cmd[arg1 arg2]
//	$cmd["arg1" "arg2"]
//	$cmd[name=value key="val"]
//	$cmd{body}
//	$cmd[args]{body}
//
// Use $$ to produce a literal $ sign.
package shortcode

import (
	"errors"
	"fmt"
	"strings"
)

// HandlerFunc is the signature for macro handler functions.
type HandlerFunc func(ctx *Context, args []string, body string) string

// PositionError wraps an error with the line and column of the macro that
// produced it. Both are 1-based.
type PositionError struct {
	Line int
	Col  int
	Err  error
}

func (e *PositionError) Error() string {
	return fmt.Sprintf("line %d, col %d: %v", e.Line, e.Col, e.Err)
}

func (e *PositionError) Unwrap() error { return e.Err }

// Context holds registered macro handlers and processes documents.
type Context struct {
	Tags map[string]HandlerFunc
	errs []error
	line int // line of the macro currently being executed (1-based)
	col  int // column of the macro currently being executed (1-based)
}

// New creates a new Context.
func New() *Context {
	return &Context{Tags: make(map[string]HandlerFunc)}
}

// AddError records an error from within a handler, annotated with the line
// and column of the macro currently being executed. The document continues
// rendering; errors are retrievable via Err or Errs after the fact.
func (c *Context) AddError(err error) {
	if err != nil {
		c.errs = append(c.errs, &PositionError{Line: c.line, Col: c.col, Err: err})
	}
}

// Err returns the first recorded error, or nil.
func (c *Context) Err() error {
	if len(c.errs) == 0 {
		return nil
	}
	return c.errs[0]
}

// Errs returns all recorded errors.
func (c *Context) Errs() []error {
	return c.errs
}

// RenderDocument renders input like Render, but clears any previous errors
// first and returns the first error (if any) alongside the output.
// Use this as the top-level entry point; use Render for recursive calls
// within handlers.
func (c *Context) RenderDocument(input string) (string, error) {
	c.errs = nil
	out := c.Render(input)
	return out, errors.Join(c.errs...)
}

// Register adds a handler for the given macro name.
func (c *Context) Register(name string, fn HandlerFunc) {
	c.Tags[name] = fn
}

// Render processes the input text, expanding macros using registered handlers.
// Unknown macros (valid syntax but unregistered name) are passed through unchanged.
// $$ produces a literal $.
func (c *Context) Render(input string) string {
	var sb strings.Builder
	i := 0
	for i < len(input) {
		if input[i] != '$' {
			sb.WriteByte(input[i])
			i++
			continue
		}
		// $$ → literal $
		if i+1 < len(input) && input[i+1] == '$' {
			sb.WriteByte('$')
			i += 2
			continue
		}
		name, args, body, end, ok := parseMacro(input, i)
		if ok {
			if fn, exists := c.Tags[name]; exists {
				c.line, c.col = lineCol(input, i)
				sb.WriteString(fn(c, args, body))
				i = end
				continue
			}
			// Valid syntax but unknown command: pass through literally.
			sb.WriteString(input[i:end])
			i = end
			continue
		}
		// Not parseable as a macro: output $ literally.
		sb.WriteByte('$')
		i++
	}
	return sb.String()
}

// parseMacro parses a macro starting at pos (which must be '$').
// Returns (name, args, body, endPos, ok).
func parseMacro(s string, pos int) (string, []string, string, int, bool) {
	i := pos + 1 // skip '$'

	// Command name: [a-zA-Z0-9_]+
	start := i
	for i < len(s) && isIdentChar(s[i]) {
		i++
	}
	if i == start {
		return "", nil, "", pos, false
	}
	name := s[start:i]

	// Optional args [...]
	var args []string
	if i < len(s) && s[i] == '[' {
		parsed, end, ok := parseArgs(s, i)
		if !ok {
			return "", nil, "", pos, false
		}
		args = parsed
		i = end
	}

	// Optional body {...}
	var body string
	if i < len(s) && s[i] == '{' {
		parsed, end, ok := parseBody(s, i)
		if !ok {
			return "", nil, "", pos, false
		}
		body = parsed
		i = end
	}

	return name, args, body, i, true
}

// parseArgs parses a bracket-delimited argument list starting at pos ('[').
// Returns (args, endPos, ok). Named args are returned as "key=value" strings.
func parseArgs(s string, pos int) ([]string, int, bool) {
	i := pos + 1 // skip '['
	var args []string

	for {
		// Skip whitespace.
		for i < len(s) && isSpace(s[i]) {
			i++
		}
		if i >= len(s) {
			return nil, pos, false
		}
		if s[i] == ']' {
			return args, i + 1, true
		}

		key, end, ok := readToken(s, i)
		if !ok {
			return nil, pos, false
		}
		i = end

		// Named arg: key=value
		if i < len(s) && s[i] == '=' {
			i++ // skip '='
			val, end2, ok2 := readToken(s, i)
			if !ok2 {
				return nil, pos, false
			}
			args = append(args, key+"="+val)
			i = end2
		} else {
			args = append(args, key)
		}
	}
}

// readToken reads a quoted or unquoted token from s starting at pos.
func readToken(s string, pos int) (string, int, bool) {
	i := pos
	if i >= len(s) {
		return "", pos, false
	}
	if s[i] == '"' {
		i++ // skip opening quote
		var sb strings.Builder
		for i < len(s) && s[i] != '"' {
			if s[i] == '\\' && i+1 < len(s) {
				i++ // skip backslash
			}
			sb.WriteByte(s[i])
			i++
		}
		if i >= len(s) {
			return "", pos, false // unterminated quote
		}
		return sb.String(), i + 1, true // skip closing quote
	}
	// Unquoted: read until whitespace, ], or =.
	start := i
	for i < len(s) && !isSpace(s[i]) && s[i] != ']' && s[i] != '=' {
		i++
	}
	if i == start {
		return "", pos, false
	}
	return s[start:i], i, true
}

// parseBody parses a brace-delimited body starting at pos ('{').
// Handles nested braces. Returns (body, endPos, ok).
func parseBody(s string, pos int) (string, int, bool) {
	i := pos + 1 // skip '{'
	depth := 1
	start := i
	for i < len(s) {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start:i], i + 1, true
			}
		}
		i++
	}
	return "", pos, false // unclosed brace
}

// ParseNamedArgs extracts key=value pairs from an args slice into a map.
// Positional args (no '=') are ignored.
func ParseNamedArgs(args []string) map[string]string {
	m := make(map[string]string)
	for _, a := range args {
		k, v, ok := strings.Cut(a, "=")
		if ok {
			m[k] = v
		}
	}
	return m
}

// lineCol returns the 1-based line and column of position pos in s.
func lineCol(s string, pos int) (line, col int) {
	line, col = 1, 1
	for i := 0; i < pos && i < len(s); i++ {
		if s[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return
}

func isIdentChar(b byte) bool {
	return b == '_' || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
