package protocol

import "testing"

func TestSimpleString_Serialize(t *testing.T) {
	s := SimpleString{Value: "OK"}
	got := string(s.Serialize())
	want := "+OK\r\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
