package jsonutil

import (
	"testing"
)

func TestDecodeEncode(t *testing.T) {
	expected := `{"bar":2,"baz":["a","b"],"foo":1}`
	s := `{"foo":1,"bar":2,"baz":["a", "b"]}`
	v, err := Decode(s)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := Encode(v)
	if err != nil {
		t.Fatal(err)
	}
	if s2 != expected {
		t.Fatalf("got %s, expected %s.", s2, expected)
	}
}
