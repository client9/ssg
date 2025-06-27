package htmlcontent

import (
	"bytes"
	"testing"
)

func TestRenderHTML(t *testing.T) {
	in := []byte("<p>test")
	want := []byte("<p>test</p>")
	out := bytes.Buffer{}

	err := Render(&out, bytes.NewReader(in), nil)
	if err != nil {
		t.Errorf("Got error in HTMLRender: %v", err)
	}
	got := out.Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("HTML want %s, got %s", want, got)
	}
}
