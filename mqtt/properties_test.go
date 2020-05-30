package mqtt

import (
	"testing"
)

func TestParseProperties(t *testing.T) {
	p := Packet{
		remainingBytes: []byte{},
	}
	a, b, c := p.parseProperties(0)
	if a != 0 {
		t.Errorf("expected a %d: %d", a, 0)
	}
	if b != 0 {
		t.Errorf("expected b %d: %d", b, 0)
	}
	if c != MALFORMED_PACKET {
		t.Errorf("expected c %d: %d", c, 0)
	}
}
