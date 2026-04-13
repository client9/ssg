package shortcode

import "testing"

func TestEntity(t *testing.T) {
	ctx := New()
	ctx.Register("ent", Entity)

	cases := []struct {
		input, want string
	}{
		{"$ent[across]", "&across;"},
		{"$ent[#1245]", "&#1245;"},
		{"$ent[amp]", "&amp;"},
		{"$ent[nbsp]", "&nbsp;"},
		{"before $ent[mdash] after", "before &mdash; after"},
	}

	for _, tc := range cases {
		got := ctx.Render(tc.input)
		if got != tc.want {
			t.Errorf("Render(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestEntityWrongArgCount(t *testing.T) {
	ctx := New()
	ctx.Register("ent", Entity)

	// Zero or multiple args return empty string.
	if got := ctx.Render("$ent[]"); got != "" {
		t.Errorf("zero args: got %q, want empty", got)
	}
	if got := ctx.Render("$ent[a b]"); got != "" {
		t.Errorf("two args: got %q, want empty", got)
	}
}
