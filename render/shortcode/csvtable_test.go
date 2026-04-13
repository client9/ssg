package shortcode

import (
	"strings"
	"testing"
)

func TestCSVTable(t *testing.T) {
	ctx := New()
	ctx.Register("csvtable", CSVTable)

	input := "$csvtable{\nName,Age,City\nAlice,25,New York\nBob,30,Boston\n}"
	got := ctx.Render(input)

	for _, want := range []string{"NAME", "Alice", "Boston", "New York"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q:\n%s", want, got)
		}
	}
}

func TestCSVTableInline(t *testing.T) {
	ctx := New()
	ctx.Register("csvtable", CSVTable)

	input := "Before\n$csvtable{Name,Score\nAlice,100\nBob,95}\nAfter"
	got := ctx.Render(input)

	if !strings.Contains(got, "Alice") || !strings.Contains(got, "100") {
		t.Errorf("unexpected output:\n%s", got)
	}
	if !strings.HasPrefix(got, "Before\n") || !strings.HasSuffix(got, "\nAfter") {
		t.Errorf("surrounding text not preserved:\n%s", got)
	}
}

// Macros in cells whose bodies contain commas must not split the CSV field.
func TestCSVTableMacroWithCommaInBody(t *testing.T) {
	ctx := New()
	ctx.Register("csvtable", CSVTable)
	ctx.Register("parens", func(_ *Context, _ []string, body string) string {
		return "(" + body + ")"
	})

	// The macro body "a,b" contains a comma — must stay in one cell.
	input := "$csvtable{Name,Value\nAlice,$parens{a,b}}"
	got := ctx.Render(input)

	if !strings.Contains(got, "(a,b)") {
		t.Errorf("macro with comma in body not rendered correctly:\n%s", got)
	}
}

// Macros should be expanded in header cells too.
func TestCSVTableMacroInHeader(t *testing.T) {
	ctx := New()
	ctx.Register("csvtable", CSVTable)
	ctx.Register("upper", func(_ *Context, _ []string, body string) string {
		return strings.ToUpper(body)
	})

	input := "$csvtable{$upper{name},Score\nAlice,100}"
	got := ctx.Render(input)

	if !strings.Contains(got, "NAME") {
		t.Errorf("macro in header not expanded:\n%s", got)
	}
}

// Macros inside a quoted CSV field.
func TestCSVTableMacroInQuotedField(t *testing.T) {
	ctx := New()
	ctx.Register("csvtable", CSVTable)
	ctx.Register("b", func(_ *Context, _ []string, body string) string {
		return "**" + body + "**"
	})

	// Quoted field contains a macro and a literal comma.
	input := `$csvtable{Name,Note` + "\n" + `Alice,"$b{hello}, world"}`
	got := ctx.Render(input)

	if !strings.Contains(got, "**hello**, world") {
		t.Errorf("macro in quoted field not rendered correctly:\n%s", got)
	}
}

func TestParseCSVWithMacros(t *testing.T) {
	cases := []struct {
		input string
		want  [][]string
	}{
		{
			"a,b,c",
			[][]string{{"a", "b", "c"}},
		},
		{
			"a,$cmd{x,y},c",
			[][]string{{"a", "$cmd{x,y}", "c"}},
		},
		{
			`"a,b",c`,
			[][]string{{"a,b", "c"}},
		},
		{
			"a\nb,c",
			[][]string{{"a"}, {"b", "c"}},
		},
		{
			`"say ""hi""",b`,
			[][]string{{"say \"hi\"", "b"}},
		},
	}

	for _, tc := range cases {
		got := parseCSVWithMacros(tc.input)
		if !rowsEqual(got, tc.want) {
			t.Errorf("input %q\n  got  %v\n  want %v", tc.input, got, tc.want)
		}
	}
}

func rowsEqual(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
