// Package email provides frontmatter parsing and serialization using
// email-style headers (RFC 822-inspired key: value pairs).
//
// This format is unusual among SSGs but offers a few advantages:
// it is human-writable without indentation rules, supports comments
// via lines starting with '#', and requires no external dependencies.
//
// Example frontmatter:
//
//	Title: My Post
//	Date: 2024-01-15
//	Tags: go, ssg, web
//	# this is a comment
//
//	<body content starts after the blank line>
package email

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ValueTransformer is a function that can inspect or transform a parsed
// key/value pair. Return io.EOF to signal the value has been fully handled
// and no further transformers should run.
type ValueTransformer func(k string, v any) (any, error)

// AsList returns a ValueTransformer that parses the named key's value
// as a CSV list rather than a plain string.
func AsList(key string) ValueTransformer {
	return func(k string, v any) (any, error) {
		if k != key {
			return v, nil
		}
		parts, err := readCSV(v.(string))
		if err != nil {
			return nil, err
		}
		return parts, io.EOF
	}
}

// ApplyTransforms runs tx in order against (k, v). If any transformer
// returns io.EOF, that value is used and iteration stops.
func ApplyTransforms(k string, v any, tx []ValueTransformer) (any, error) {
	for _, t := range tx {
		nv, err := t(k, v)
		if errors.Is(err, io.EOF) {
			return nv, nil
		}
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

// Set writes k/v into kv, applying any transformers first.
// Supports dotted keys (e.g. "author.name") to create nested maps.
func Set(kv map[string]any, k string, v any, tx []ValueTransformer) error {
	v, err := ApplyTransforms(k, v, tx)
	if err != nil {
		return err
	}

	parts := strings.Split(strings.TrimSpace(k), ".")
	base := parts[0 : len(parts)-1]
	key := parts[len(parts)-1]

	next := kv
	for _, b := range base {
		if _, ok := next[b]; !ok {
			next[b] = make(map[string]any)
		}
		tmp, ok := next[b].(map[string]any)
		if !ok {
			return fmt.Errorf("key %q isnt a map", k)
		}
		next = tmp
	}
	if _, ok := next[key]; !ok {
		next[key] = v
		return nil
	}
	return fmt.Errorf("duplicate key %q", k)
}

func replaceNewlines(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}

func readCSV(row string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(row))
	return r.Read()
}

func writeCSV(parts []string) ([]byte, error) {
	for i, s := range parts {
		parts[i] = replaceNewlines(s)
	}
	out := &bytes.Buffer{}
	w := csv.NewWriter(out)
	_ = w.Write(parts)
	w.Flush()
	return out.Bytes(), nil
}

// Unmarshal parses email-style headers from src into out.
// Lines starting with '#' are treated as comments and ignored.
// A line indented with a space or tab continues the previous value.
func Unmarshal(src []byte, out map[string]any, tx ...ValueTransformer) error {
	scan := bufio.NewScanner(bytes.NewReader(src))
	scan.Split(bufio.ScanLines)
	key := ""
	prev := ""
	for scan.Scan() {
		text := scan.Text()
		if len(text) == 0 || text[0] == '#' {
			if key != "" {
				if err := Set(out, key, prev, tx); err != nil {
					return err
				}
				key = ""
				prev = ""
			}
			continue
		}

		if text[0] == ' ' || text[0] == '\t' {
			if key == "" {
				return fmt.Errorf("invalid line starting with whitespace")
			}
			prev += " " + strings.TrimSpace(text)
			continue
		}
		if key != "" {
			if err := Set(out, key, prev, tx); err != nil {
				return err
			}
		}
		k, v, ok := strings.Cut(text, ":")
		if !ok {
			return fmt.Errorf("malformed header line: " + text)
		}
		key = k
		prev = strings.TrimSpace(v)
	}
	if err := scan.Err(); err != nil {
		return err
	}
	if key != "" {
		if err := Set(out, key, prev, tx); err != nil {
			return err
		}
	}
	return nil
}

func appendKey(out []byte, prefix, key string) []byte {
	out = append(out, []byte(prefix)...)
	out = append(out, []byte(key)...)
	out = append(out, byte(':'), byte(' '))
	return out
}

// Marshal serializes a map[string]any into email-header style output.
// Nested maps are written with dotted key prefixes (e.g. "author.name: foo").
func Marshal(data map[string]any) ([]byte, error) {
	out := make([]byte, 0, 1024)
	return writeEmailMeta(out, "", data)
}

func writeEmailMeta(out []byte, prefix string, data map[string]any) ([]byte, error) {
	keys := slices.Sorted(maps.Keys(data))
	var current, next []string
	for _, k := range keys {
		if _, ok := data[k].(map[string]any); ok {
			next = append(next, k)
		} else {
			current = append(current, k)
		}
	}

	for _, k := range current {
		v := data[k]
		switch val := v.(type) {
		case []any:
			tmp := make([]string, len(val))
			for i, item := range val {
				tmp[i] = fmt.Sprintf("%v", item)
			}
			out = appendKey(out, prefix, k)
			row, err := writeCSV(tmp)
			if err != nil {
				return nil, err
			}
			out = append(out, row...)
		case []string:
			out = appendKey(out, prefix, k)
			row, err := writeCSV(val)
			if err != nil {
				return nil, err
			}
			out = append(out, row...)
		case string:
			out = appendKey(out, prefix, k)
			row, err := writeCSV([]string{val})
			if err != nil {
				return nil, err
			}
			out = append(out, row...)
		case *time.Time:
			out = appendKey(out, prefix, k)
			out = append(out, []byte(val.String())...)
			out = append(out, '\n')
		case time.Time:
			out = appendKey(out, prefix, k)
			out = append(out, []byte(val.String())...)
			out = append(out, '\n')
		case float32:
			out = appendKey(out, prefix, k)
			out = strconv.AppendFloat(out, float64(val), 'g', -1, 32)
			out = append(out, '\n')
		case float64:
			out = appendKey(out, prefix, k)
			out = strconv.AppendFloat(out, val, 'g', -1, 64)
			out = append(out, '\n')
		case bool:
			out = appendKey(out, prefix, k)
			out = strconv.AppendBool(out, val)
			out = append(out, '\n')
		case int:
			out = appendKey(out, prefix, k)
			out = strconv.AppendInt(out, int64(val), 10)
			out = append(out, '\n')
		case int64:
			out = appendKey(out, prefix, k)
			out = strconv.AppendInt(out, val, 10)
			out = append(out, '\n')
		case uint:
			out = appendKey(out, prefix, k)
			out = strconv.AppendUint(out, uint64(val), 10)
			out = append(out, '\n')
		case uint64:
			out = appendKey(out, prefix, k)
			out = strconv.AppendUint(out, val, 10)
			out = append(out, '\n')
		default:
			return nil, fmt.Errorf("unknown type %T with value %v", v, v)
		}
	}
	if len(current) > 0 {
		out = append(out, '\n')
	}

	for _, k := range next {
		v := data[k].(map[string]any)
		tout, err := writeEmailMeta(out, k+".", v)
		if err != nil {
			return out, err
		}
		out = tout
	}

	return out, nil
}

// Loader is the default MetaLoader for email-style frontmatter with no
// value transformers. Use NewLoader to apply transformers such as AsList.
var Loader = NewLoader()

// NewLoader returns a MetaLoader for email-style frontmatter, applying tx in
// order to each parsed key/value pair.
//
//	email.NewLoader(email.AsList("Tags"))
func NewLoader(tx ...ValueTransformer) func([]byte) (map[string]any, []byte, error) {
	return func(raw []byte) (map[string]any, []byte, error) {
		head, body := split(raw)
		if head == nil {
			return map[string]any{}, body, nil
		}
		meta := map[string]any{}
		if err := Unmarshal(head, meta, tx...); err != nil {
			return nil, nil, fmt.Errorf("unable to parse metadata: %v", err)
		}
		return meta, body, nil
	}
}

func split(raw []byte) (head, body []byte) {
	head, body, found := bytes.Cut(raw, []byte("\n\n\n"))
	if !found {
		return nil, raw
	}
	return head, body
}
