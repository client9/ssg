package funcs

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		in   string
		n    int
		want string
	}{
		{"hello", 10, "hello"},         // shorter than limit
		{"hello", 5, "hello"},          // exact length
		{"hello world", 8, "hello w…"}, // cut with ellipsis
		{"hello", 1, "…"},              // n=1 is just ellipsis
		{"hello", 0, ""},               // n=0 is empty
		{"héllo", 4, "hél…"},           // rune-aware
	}
	for _, tt := range tests {
		got := truncate(tt.in, tt.n)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.in, tt.n, got, tt.want)
		}
	}
}
