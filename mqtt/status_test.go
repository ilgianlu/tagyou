package mqtt

import "testing"

func TestCleanSession(t *testing.T) {
	var c ConnStatus
	c.connectFlags = 0x02
	if !c.cleanSession() {
		t.Error("expected clean session true")
	}
	c.connectFlags = 0x00
	if c.cleanSession() {
		t.Error("expected clean session false")
	}
}
