package stdfuncs_test

import (
	"fmt"
	"html/template"
	"time"

	funcs "github.com/client9/ssg/tmpl/stdfuncs"
)

// ---- math ----

func ExampleClamp() {
	within, _ := funcs.Clamp(5, 1, 10)
	below, _ := funcs.Clamp(0, 1, 10)
	above, _ := funcs.Clamp(15, 1, 10)
	fmt.Println(within, below, above)
	// Output:
	// 5 1 10
}

func ExamplePow() {
	v, _ := funcs.Pow(2, 10)
	fmt.Println(v)
	// Output:
	// 1024
}

func ExamplePow_sqrt() {
	v, _ := funcs.Pow(9, 0.5)
	fmt.Println(v)
	// Output:
	// 3
}

func ExampleModBool() {
	even, _ := funcs.ModBool(4, 2)
	odd, _ := funcs.ModBool(5, 2)
	fmt.Println(even, odd)
	// Output:
	// true false
}

func ExampleMin() {
	v, _ := funcs.Min(3, 1, 4, 1, 5, 9)
	fmt.Println(v)
	// Output:
	// 1
}

func ExampleMin_slice() {
	v, _ := funcs.Min([]int{7, 2, 8})
	fmt.Println(v)
	// Output:
	// 2
}

func ExampleMax() {
	v, _ := funcs.Max(3, 1, 4, 1, 5, 9)
	fmt.Println(v)
	// Output:
	// 9
}

// ---- collections: constructors ----

func ExampleList() {
	s := funcs.List("a", "b", "c")
	fmt.Println(s)
	// Output:
	// [a b c]
}

func ExampleDict() {
	m, _ := funcs.Dict("name", "Alice", "age", 30)
	fmt.Println(m["name"], m["age"])
	// Output:
	// Alice 30
}

func ExampleSeq() {
	s, _ := funcs.Seq(5)
	fmt.Println(s)
	// Output:
	// [1 2 3 4 5]
}

func ExampleSeq_range() {
	s, _ := funcs.Seq(3, 7)
	fmt.Println(s)
	// Output:
	// [3 4 5 6 7]
}

func ExampleSeq_step() {
	s, _ := funcs.Seq(1, 9, 2)
	fmt.Println(s)
	// Output:
	// [1 3 5 7 9]
}

// ---- collections: sequence access ----

func ExampleFirst() {
	v, _ := funcs.First([]string{"a", "b", "c"})
	fmt.Println(v)
	// Output:
	// a
}

func ExampleFirst_string() {
	v, _ := funcs.First("café")
	fmt.Println(v)
	// Output:
	// c
}

func ExampleLast() {
	v, _ := funcs.Last([]string{"a", "b", "c"})
	fmt.Println(v)
	// Output:
	// c
}

func ExampleLast_string() {
	v, _ := funcs.Last("café")
	fmt.Println(v)
	// Output:
	// é
}

func ExampleTake() {
	v, _ := funcs.Take([]int{1, 2, 3, 4, 5}, 3)
	fmt.Println(v)
	// Output:
	// [1 2 3]
}

func ExampleTake_negative() {
	v, _ := funcs.Take([]int{1, 2, 3, 4, 5}, -2)
	fmt.Println(v)
	// Output:
	// [4 5]
}

func ExampleTake_string() {
	v, _ := funcs.Take("日本語", 2)
	fmt.Println(v)
	// Output:
	// 日本
}

func ExampleDrop() {
	v, _ := funcs.Drop([]int{1, 2, 3, 4, 5}, 2)
	fmt.Println(v)
	// Output:
	// [3 4 5]
}

func ExampleDrop_negative() {
	v, _ := funcs.Drop([]int{1, 2, 3, 4, 5}, -2)
	fmt.Println(v)
	// Output:
	// [1 2 3]
}

func ExampleDrop_string() {
	v, _ := funcs.Drop("hello", 2)
	fmt.Println(v)
	// Output:
	// llo
}

// ---- collections: sequence transformation ----

func ExampleReverse() {
	v, _ := funcs.Reverse([]int{1, 2, 3})
	fmt.Println(v)
	// Output:
	// [3 2 1]
}

func ExampleCompact() {
	v, _ := funcs.Compact([]any{1, 1, 2, 3, 3, 1})
	fmt.Println(v)
	// Output:
	// [1 2 3 1]
}

func ExampleConcat() {
	v, _ := funcs.Concat([]int{1, 2}, []int{3, 4})
	fmt.Println(v)
	// Output:
	// [1 2 3 4]
}

func ExampleSort() {
	v, _ := funcs.Sort([]string{"banana", "apple", "cherry"})
	fmt.Println(v)
	// Output:
	// [apple banana cherry]
}

func ExampleSort_numeric() {
	v, _ := funcs.Sort([]int{10, 2, 30, 5})
	fmt.Println(v)
	// Output:
	// [2 5 10 30]
}

func ExampleSort_time() {
	t1 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	v, _ := funcs.Sort([]time.Time{t1, t2})
	fmt.Println(v.([]any)[0].(time.Time).Format("2006-01-02"))
	// Output:
	// 2023-12-31
}

func ExampleSortNum() {
	v, _ := funcs.SortNum([]string{"10", "9", "2"})
	fmt.Println(v)
	// Output:
	// [2 9 10]
}

func ExampleWhere() {
	pages := []any{
		map[string]any{"Title": "Post A", "Draft": false},
		map[string]any{"Title": "Post B", "Draft": true},
		map[string]any{"Title": "Post C", "Draft": false},
	}
	v, _ := funcs.Where(pages, "Draft", false)
	fmt.Println(len(v.([]any)))
	// Output:
	// 2
}

// ---- collections: map operations ----

func ExampleKeys() {
	m := map[string]any{"b": 2, "a": 1, "c": 3}
	k, _ := funcs.Keys(m)
	fmt.Println(k)
	// Output:
	// [a b c]
}

func ExampleValues() {
	m := map[string]any{"b": 2, "a": 1}
	v, _ := funcs.Values(m)
	fmt.Println(v)
	// Output:
	// [1 2]
}

func ExampleMergeMaps() {
	a := map[string]any{"x": 1, "y": 2}
	b := map[string]any{"y": 99, "z": 3}
	m, _ := funcs.MergeMaps(a, b)
	fmt.Println(m["x"], m["y"], m["z"])
	// Output:
	// 1 99 3
}

// ---- collections: general ----

func ExampleIn_slice() {
	ok, _ := funcs.In([]string{"a", "b", "c"}, "b")
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleIn_map() {
	ok, _ := funcs.In(map[string]any{"x": 1}, "x")
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleIn_string() {
	ok, _ := funcs.In("hello world", "world")
	fmt.Println(ok)
	// Output:
	// true
}

func ExampleDefault() {
	fmt.Println(funcs.Default("anon", ""))
	fmt.Println(funcs.Default("anon", "Alice"))
	// Output:
	// anon
	// Alice
}

func ExampleCond() {
	fmt.Println(funcs.Cond(true, "yes", "no"))
	fmt.Println(funcs.Cond(false, "yes", "no"))
	fmt.Println(funcs.Cond(0, "yes", "no"))
	// Output:
	// yes
	// no
	// no
}

// ---- safe types ----

func ExampleSafeHTML() {
	v, _ := funcs.SafeHTML("<b>bold</b>")
	fmt.Println(v)
	// Output:
	// <b>bold</b>
}

func ExampleSafeCSS() {
	v, _ := funcs.SafeCSS("color: red")
	fmt.Println(v)
	// Output:
	// color: red
}

func ExampleSafeURL() {
	v, _ := funcs.SafeURL("https://example.com/path?q=1")
	fmt.Println(v)
	// Output:
	// https://example.com/path?q=1
}

func ExampleSafeJS() {
	v, _ := funcs.SafeJS("alert('hi')")
	fmt.Println(v)
	// Output:
	// alert('hi')
}

func ExampleSafeHTMLAttr() {
	v, _ := funcs.SafeHTMLAttr(`class="hero"`)
	fmt.Println(v)
	// Output:
	// class="hero"
}

func ExampleSafeJSStr() {
	v, _ := funcs.SafeJSStr(`hello\nworld`)
	fmt.Println(v)
	// Output:
	// hello\nworld
}

// ---- strings ----

func ExampleCapitalize() {
	fmt.Println(funcs.Capitalize("hello world"))
	fmt.Println(funcs.Capitalize("NASA"))
	// Output:
	// Hello world
	// Nasa
}

func ExampleLenRunes() {
	fmt.Println(funcs.LenRunes("café"))
	fmt.Println(funcs.LenRunes("日本語"))
	// Output:
	// 4
	// 3
}

func ExampleReplace() {
	fmt.Println(funcs.Replace("aabbaa", "a", "x"))
	// Output:
	// xabbaa
}

func ExampleReplace_count() {
	fmt.Println(funcs.Replace("aabbaa", "a", "x", -1))
	// Output:
	// xxbbxx
}

// ---- encoding ----

func ExampleJsonify() {
	v, _ := funcs.Jsonify(map[string]any{"name": "Alice", "age": 30})
	fmt.Println(v)
	// Output:
	// {"age":30,"name":"Alice"}
}

func ExampleJsonify_slice() {
	v, _ := funcs.Jsonify([]string{"a", "b", "c"})
	fmt.Println(v)
	// Output:
	// ["a","b","c"]
}

// ---- cast ----

func ExampleToInt() {
	v, _ := funcs.ToInt("42")
	fmt.Println(v)
	// Output:
	// 42
}

func ExampleToInt_float() {
	v, _ := funcs.ToInt(3.9)
	fmt.Println(v)
	// Output:
	// 3
}

func ExampleToFloat() {
	v, _ := funcs.ToFloat("3.14")
	fmt.Println(v)
	// Output:
	// 3.14
}

// ---- time ----

func ExampleNow() {
	t := funcs.Now()
	_ = t
	// Returns current local time as time.Time.
	// In templates: {{now.Year}}, {{now.Format "2006-01-02"}}
}

func ExampleParseTime() {
	t, _ := funcs.ParseTime("2006-01-02", "2024-03-15")
	fmt.Println(t.Format("January 2, 2006"))
	// Output:
	// March 15, 2024
}

// ---- FuncMap / Merge ----

func ExampleFuncMap() {
	fm := funcs.FuncMap()
	t := template.Must(template.New("").Funcs(fm).Parse(`{{upper "hello"}}`))
	_ = t.Execute(nil, nil)
	// (FuncMap registers all template functions; use with template.New().Funcs())
}

func ExampleMerge() {
	custom := template.FuncMap{
		"greet": func(name string) string { return "Hello, " + name + "!" },
	}
	fm := funcs.Merge(funcs.FuncMap(), custom)
	t := template.Must(template.New("").Funcs(fm).Parse(`{{greet "World"}}`))
	_ = t.Execute(nil, nil)
	// (Merge combines FuncMaps; later maps win on key collision)
}
