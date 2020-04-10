package mqtt

import "testing"

func TestCleanSession(t *testing.T) {
	var c Connection
	c.connectFlags = 0x02
	if !c.cleanStart() {
		t.Error("expected clean session true")
	}
	c.connectFlags = 0x00
	if c.cleanStart() {
		t.Error("expected clean session false")
	}
}
