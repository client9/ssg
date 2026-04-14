package htmlclean

import (
	"testing"
)

func TestRenderHTML(t *testing.T) {
	in := []byte("<p>test")
	want := []byte("<p>test</p>")

	out, err := Render.Run(nil, nil, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.([]byte)
	if string(got) != string(want) {
		t.Errorf("HTML want %s, got %s", want, got)
	}
}
