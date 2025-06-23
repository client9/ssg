package ssg

import (
	"bytes"
	"testing"
)

func TestIdenity(t *testing.T) {
	in := []byte("123")
	out := bytes.Buffer{}
	err := Identity(&out, bytes.NewReader(in), nil)
	if err != nil {
		t.Errorf("Got error in HTMLRender: %v", err)
	}
	want := in
	got := out.Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("Identity want %s, got %s", want, got)
	}
}

func TestRenderHTML(t *testing.T) {
	in := []byte("<p>test")
	want := []byte("<p>test</p>")
	out := bytes.Buffer{}

	err := HTMLRender(&out, bytes.NewReader(in), nil)
	if err != nil {
		t.Errorf("Got error in HTMLRender: %v", err)
	}
	got := out.Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("HTML want %s, got %s", want, got)
	}
}

func TestRenderTemplateMacro(t *testing.T) {

	render := NewTemplateMacro(nil)

	in := []byte(`<p>{{ printf "123" }}</p>`)
	want := []byte("<p>123</p>")
	out := bytes.Buffer{}

	err := render(&out, bytes.NewReader(in), nil)
	if err != nil {
		t.Errorf("Got error in MacroRenderer: %v", err)
	}
	got := out.Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("MacroRenderer want %s, got %s", want, got)
	}
}

func TestMultiRender(t *testing.T) {

	in := []byte(`<p >{{ printf "123" }}</p >`)
	want := []byte("<p>123</p>")
	out := bytes.Buffer{}

	pipeline := []Renderer{
		Identity,
		NewTemplateMacro(nil),
		HTMLRender,
		ToBytes(&out),
	}

	err := MultiRender(pipeline, in, nil)

	if err != nil {
		t.Errorf("Got error in MultiRenderer: %v", err)
	}
	got := out.Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("MultiRenderer want %s, got %s", want, got)
	}
}
