package packet

import (
	"testing"
)

func TestParseProperties(t *testing.T) {
	p := Packet{
		remainingBytes: []byte{},
	}
	a, c := p.parseProperties(0)
	if a != 0 {
		t.Errorf("expected a %d: %d", a, 0)
	}
	if c != MALFORMED_PACKET {
		t.Errorf("expected c %d: %d", c, MALFORMED_PACKET)
	}
	p0 := Packet{
		remainingBytes: []byte{
			7,
			SESSION_EXPIRY_INTERVAL, 0, 0, 0, 0xA,
			SHARED_SUBSCRIPTION_AVAILABLE, 0},
	}
	a0, c0 := p0.parseProperties(0)
	if a0 != 8 {
		t.Errorf("expected a0 %d: %d", a0, 8)
	}
	if c0 != 0 {
		t.Errorf("expected c0 %d: %d", c0, 0)
	}
	if p0.properties[SESSION_EXPIRY_INTERVAL].position != 2 {
		t.Errorf("expected exp int pos %d: %d", p0.properties[SESSION_EXPIRY_INTERVAL].position, 2)
	}
	if p0.properties[SESSION_EXPIRY_INTERVAL].length != 4 {
		t.Errorf("expected exp int pos %d: %d", p0.properties[SESSION_EXPIRY_INTERVAL].length, 4)
	}
	if p0.properties[SHARED_SUBSCRIPTION_AVAILABLE].position != 7 {
		t.Errorf("expected exp int pos %d: %d", p0.properties[SHARED_SUBSCRIPTION_AVAILABLE].position, 7)
	}
	if p0.properties[SHARED_SUBSCRIPTION_AVAILABLE].length != 1 {
		t.Errorf("expected exp int pos %d: %d", p0.properties[SHARED_SUBSCRIPTION_AVAILABLE].length, 1)
	}
	seir := p0.getPropertyRaw(SESSION_EXPIRY_INTERVAL)
	if seir[3] != 10 {
		t.Errorf("expected ses exp int raw %d: %d", seir[3], 10)
	}
	sei := p0.SessionExpiryInterval()
	if sei != 10 {
		t.Errorf("expected ses exp int %d: %d", seir[3], 10)
	}
}
