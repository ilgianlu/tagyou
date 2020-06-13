package mqtt

import (
	"testing"
)

func TestConnack(t *testing.T) {
	p := Connack(false, 0, MQTT_V5)
	if p.PacketType() != PACKET_TYPE_CONNACK {
		t.Errorf("Connack expected packet type %d, got %d", PACKET_TYPE_CONNACK, p.PacketType())
	}
	if p.Flags() != 0 {
		t.Errorf("Connack expected flags %d, got %d", 0, p.Flags())
	}
	if p.remainingLength != 3 {
		t.Errorf("Connack remaining length %d, got %d", 3, p.remainingLength)
	}
	p = Connack(false, 0, MQTT_V3_11)
	if p.remainingLength != 2 {
		t.Errorf("Connack remaining length %d, got %d", 2, p.remainingLength)
	}
}
