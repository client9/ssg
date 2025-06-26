package ssg

import (
	"time"
	"testing"
)

func TestEmailMeta(t *testing.T) {

	msg := []byte(`
first: 1
# comment

second: 2
`)
	out, err := EmailMeta(msg)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("Expected 2 keys, got %d %v", len(out), out)
	}
}
func TestEmailMetaContinued(t *testing.T) {

	msg := []byte(`
first: This is 
 a line
second: 2
`)
	out, err := EmailMeta(msg)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("Expected 2 keys, got %d %v", len(out), out)
	}
	val := out["first"].(string)
	want := "This is a line"
	if val != want {
		t.Errorf("Expected %q, got %q", want, val)
	}
}

func TestEmailWrite1(t *testing.T) {

	data := make(map[string]any)
	data["afloat"] = 1.125
	data["aint"] = 124
	data["aint64"] = int64(125)
	data["auint64"] = uint64(999)
	data["abool"] = true
	data["astring"] = "hello world"
	data["tag"] = []string{"apple", "banana"}
	data["atime"] = time.Date(2025,1,2,1,2,3,0,time.UTC)
	want := `abool: true
afloat: 1.125
aint: 124
aint64: 125
astring: hello world
atime: 2025-01-02 01:02:03 +0000 UTC
auint64: 999
tag: apple
tag: banana
`

	out, err := EmailMarshal(data)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	s := string(out)
	if s != want {
		t.Errorf("Got\n %s", s)
	}
}
func TestEmailWrite2(t *testing.T) {

}
