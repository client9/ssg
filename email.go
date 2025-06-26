package ssg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
)

type KVTransformer func(kv map[string]any, k, v string) error

func KVForKey(key string, next KVTransformer) KVTransformer {
	return func(kv map[string]any, k, v string) error {
		if k == key {
			return next(kv, k, v)
		}
		return nil
	}
}

func KVOverwrite(kv map[string]any, k, v string) error {
	kv[k] = v
	return io.EOF
}

func KVDupError(kv map[string]any, k, v string) error {
	_, ok := kv[k]
	if !ok {
		kv[k] = v
		return io.EOF
	}
	return fmt.Errorf("duplicate key %q", k)
}

func KVStringAllList(kv map[string]any, k, v string) error {
	val, ok := kv[k]
	if !ok {
		kv[k] = []string{v}
		return io.EOF
	}
	list := val.([]string)
	list = append(list, v)
	kv[k] = list
	return io.EOF
}

func KVStringList(kv map[string]any, k, v string) error {
	val, ok := kv[k]
	if !ok {
		kv[k] = []string{v}
		return io.EOF
	}
	list := val.([]string)
	list = append(list, v)
	kv[k] = list
	return io.EOF
}

func ApplyTx(kv map[string]any, k, v string, tx []KVTransformer) error {
	for _, t := range tx {
		err := t(kv, k, v)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func EmailMeta(src []byte, tx ...KVTransformer) (map[string]any, error) {
	if len(tx) == 0 {
		tx = append(tx, KVOverwrite)
	}
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
				if err := ApplyTx(out, key, prev, tx); err != nil {
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
			if err := ApplyTx(out, key, prev, tx); err != nil {
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
		if err := ApplyTx(out, key, prev, tx); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// write out data is more straightfoward.
// The only ootion is really what to do with []string values
//  One header per value "foo:bar1\nfoo:bar2"
//  CSV "foo: bar1,bar2" .. might need quoting issues.

func appendKey(out []byte, key string) []byte {
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
	for _, k := range keys {
		v := data[k]

		switch val := v.(type) {
		case []any:
			// yaml does this
			tmp := make([]string, len(val))
			for i, v := range val {
				tmp[i] = fmt.Sprintf("%v", v)
			}
			out = appendKey(out, k)
			out = append(out, []byte(strings.Join(tmp, ", "))...)
			out = append(out, byte('\n'))
		case *time.Time:
			out = appendKey(out, k)
			out = append(out, []byte(val.String())...)
			out = append(out, byte('\n'))
		case time.Time:
			out = appendKey(out, k)
			out = append(out, []byte(val.String())...)
			out = append(out, byte('\n'))
		case map[string]any:
			out, err := writeEmailMeta(out, k+"-", val)
			if err != nil {
				return out, err
			}
		case string:
			out = appendKey(out, k)
			out = append(out, []byte(val)...)
			out = append(out, byte('\n'))
		case []string:
			for _, s := range val {
				out = appendKey(out, k)
				out = append(out, []byte(s)...)
				out = append(out, byte('\n'))
			}
		case float32:
			out = appendKey(out, k)
			out = strconv.AppendFloat(out, float64(val), 'g', -1, 32)
			out = append(out, byte('\n'))
		case float64:
			out = appendKey(out, k)
			out = strconv.AppendFloat(out, val, 'g', -1, 64)
			out = append(out, byte('\n'))
		case bool:
			out = appendKey(out, k)
			out = strconv.AppendBool(out, val)
			out = append(out, byte('\n'))
		case int:
			out = appendKey(out, k)
			out = strconv.AppendInt(out, int64(val), 10)
			out = append(out, byte('\n'))
		case int64:
			out = appendKey(out, k)
			out = strconv.AppendInt(out, val, 10)
			out = append(out, byte('\n'))
		case uint:
			out = appendKey(out, k)
			out = strconv.AppendUint(out, uint64(val), 10)
			out = append(out, byte('\n'))
		case uint64:
			out = appendKey(out, k)
			out = strconv.AppendUint(out, val, 10)
			out = append(out, byte('\n'))
		default:
			return nil, fmt.Errorf("unknown type %T with value %v", v, v)
		}
	}
	return out, nil
}
