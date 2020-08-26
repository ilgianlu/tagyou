package packet

import (
	"testing"
)

func TestPublish(t *testing.T) {
	p := Publish(4, 0, false, "topic", 0, []byte{0, 1, 2})
	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("Publish expected remaingLength %d, received %d", len(p.remainingBytes), p.remainingLength)
	}
	p = Publish(4, 1, false, "topic", 123, []byte{0, 1, 2})
	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("Publish expected remaingLength %d, received %d", len(p.remainingBytes), p.remainingLength)
	}
	if p.QoS() != 1 {
		t.Errorf("Publish expected qos 1, received %d", p.QoS())
	}
}
