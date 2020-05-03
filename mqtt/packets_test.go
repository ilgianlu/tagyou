package mqtt

import (
	"testing"
)

func TestPublish(t *testing.T) {
	p := Publish(0, false, "topic", []byte{0, 1, 2})
	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("Publish expected remaingLength %d, received %d", len(p.remainingBytes), p.remainingLength)
	}
}

func TestConnack(t *testing.T) {
	p := Connack(false, 0)
	if p.PacketType() != PACKET_TYPE_CONNACK {
		t.Errorf("Connack expected packet type %d, got %d", PACKET_TYPE_CONNACK, p.PacketType())
	}
	if p.Flags() != 0 {
		t.Errorf("Connack expected flags %d, got %d", 0, p.Flags())
	}
	if p.remainingLength != 2 {
		t.Errorf("Connack remaining length %d, got %d", 2, p.remainingLength)
	}
}
