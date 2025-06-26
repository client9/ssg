package ssg

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

type ValueTransformer func(k string, v any) (any, error)

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

func ApplyTransforms(k string, v any, tx []ValueTransformer) (any, error) {
	for _, t := range tx {
		v, err := t(k, v)
		if errors.Is(err, io.EOF) {
			// we handled it, done.
			return v, nil
		}
		if err != nil {
			// something didn't work
			return nil, err
		}
		// didn't process it, keep going
	}

	// no change
	return v, nil
}

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
			next = make(map[string]any)
			kv[b] = next
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
	// one row, one shot!
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

func EmailMeta(src []byte, tx ...ValueTransformer) (map[string]any, error) {
	out := make(map[string]any)
	scan := bufio.NewScanner(bytes.NewReader(src))
	scan.Split(bufio.ScanLines)
	key := ""
	prev := ""
	for scan.Scan() {
		text := scan.Text()
		// empty line or comment line
		if len(text) == 0 || text[0] == '#' {
			if key != "" {
				if err := Set(out, key, prev, tx); err != nil {
					return nil, err
				}
				key = ""
				prev = ""
			}
			continue
		}

		if text[0] == ' ' || text[0] == '\t' {
			if key == "" {
				return nil, fmt.Errorf("invalid line starting with whitespace")
			}
			prev += " " + strings.TrimSpace(text)
			continue
		}
		// this is a normal line.
		if key != "" {
			if err := Set(out, key, prev, tx); err != nil {
				return nil, err
			}
		}
		k, v, ok := strings.Cut(text, ":")
		if !ok {
			return nil, fmt.Errorf("malformed header line: " + text)
		}
		key = k
		prev = strings.TrimSpace(v)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	if key != "" {
		if err := Set(out, key, prev, tx); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func appendKey(out []byte, prefix, key string) []byte {
	out = append(out, []byte(prefix)...)
	out = append(out, []byte(key)...)
	out = append(out, byte(':'), byte(' '))
	return out
}

func EmailMarshal(data map[string]any) ([]byte, error) {
	out := make([]byte, 0, 1024)
	return writeEmailMeta(out, "", data)
}

func writeEmailMeta(out []byte, prefix string, data map[string]any) ([]byte, error) {
	keys := slices.Sorted(maps.Keys(data))
	current := []string{}
	next := []string{}
	for _, k := range keys {
		v := data[k]
		switch v.(type) {
		case map[string]any:
			next = append(next, k)
		default:
			current = append(current, k)
		}
	}

	for _, k := range current {
		v := data[k]

		switch val := v.(type) {
		case []any:
			// yaml uses []any
			tmp := make([]string, len(val))
			for i, v := range val {
				tmp[i] = fmt.Sprintf("%v", v)
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
			list := []string{val}
			row, err := writeCSV(list)
			if err != nil {
				return nil, err
			}
			out = append(out, row...)
		case *time.Time:
			out = appendKey(out, prefix, k)
			out = append(out, []byte(val.String())...)
			out = append(out, byte('\n'))
		case time.Time:
			out = appendKey(out, prefix, k)
			out = append(out, []byte(val.String())...)
			out = append(out, byte('\n'))
		case float32:
			out = appendKey(out, prefix, k)
			out = strconv.AppendFloat(out, float64(val), 'g', -1, 32)
			out = append(out, byte('\n'))
		case float64:
			out = appendKey(out, prefix, k)
			out = strconv.AppendFloat(out, val, 'g', -1, 64)
			out = append(out, byte('\n'))
		case bool:
			out = appendKey(out, prefix, k)
			out = strconv.AppendBool(out, val)
			out = append(out, byte('\n'))
		case int:
			out = appendKey(out, prefix, k)
			out = strconv.AppendInt(out, int64(val), 10)
			out = append(out, byte('\n'))
		case int64:
			out = appendKey(out, prefix, k)
			out = strconv.AppendInt(out, val, 10)
			out = append(out, byte('\n'))
		case uint:
			out = appendKey(out, prefix, k)
			out = strconv.AppendUint(out, uint64(val), 10)
			out = append(out, byte('\n'))
		case uint64:
			out = appendKey(out, prefix, k)
			out = strconv.AppendUint(out, val, 10)
			out = append(out, byte('\n'))
		default:
			return nil, fmt.Errorf("unknown type %T with value %v", v, v)
		}
	}
	if len(current) > 0 {
		out = append(out, byte('\n'))
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
