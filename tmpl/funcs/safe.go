package funcs

import (
	"fmt"
	"html/template"
	"net/url"
)

func safeFuncMap() template.FuncMap {
	return template.FuncMap{
		"safeCSS":       SafeCSS,
		"safeHTML":      SafeHTML,
		"safeHTMLAttr":  SafeHTMLAttr,
		"safeJS":        SafeJS,
		"safeJSStr":     SafeJSStr,
		"safeURL":       SafeURL,
		"urlEncode":     url.QueryEscape,
		"urlPathEscape": url.PathEscape,
	}
}

// safeString coerces any value to a plain string for wrapping in a safe type.
// html/template typed values (e.g. template.HTML) are unwrapped without re-encoding.
func safeString(s any) (string, error) {
	switch v := s.(type) {
	case string:
		return v, nil
	case template.HTML:
		return string(v), nil
	case template.CSS:
		return string(v), nil
	case template.HTMLAttr:
		return string(v), nil
	case template.JS:
		return string(v), nil
	case template.JSStr:
		return string(v), nil
	case template.URL:
		return string(v), nil
	case []byte:
		return string(v), nil
	case nil:
		return "", fmt.Errorf("safe: nil input")
	default:
		return fmt.Sprint(v), nil
	}
}

// SafeCSS converts s to template.CSS, marking it safe for use in style attributes
// and <style> blocks without escaping.
func SafeCSS(s any) (template.CSS, error) {
	str, err := safeString(s)
	return template.CSS(str), err
}

// SafeHTML converts s to template.HTML, marking it safe to render as raw HTML
// without escaping. Use only with trusted content.
func SafeHTML(s any) (template.HTML, error) {
	str, err := safeString(s)
	return template.HTML(str), err
}

// SafeHTMLAttr converts s to template.HTMLAttr, marking it safe for use as an
// HTML attribute (name and value pair) without escaping.
func SafeHTMLAttr(s any) (template.HTMLAttr, error) {
	str, err := safeString(s)
	return template.HTMLAttr(str), err
}

// SafeJS converts s to template.JS, marking it safe for use inside <script>
// blocks without escaping.
func SafeJS(s any) (template.JS, error) {
	str, err := safeString(s)
	return template.JS(str), err
}

// SafeJSStr converts s to template.JSStr, marking it safe for interpolation
// inside a JavaScript string literal without escaping.
func SafeJSStr(s any) (template.JSStr, error) {
	str, err := safeString(s)
	return template.JSStr(str), err
}

// SafeURL converts s to template.URL, marking it safe for use in URL attributes
// (href, src, action, etc.) without escaping.
func SafeURL(s any) (template.URL, error) {
	str, err := safeString(s)
	return template.URL(str), err
}
