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

func TestStart(t *testing.T) {
	b := []byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46, 46, 0, 4, 117, 115, 101, 114, 0, 5, 112, 108, 117, 116, 111}
	p, err := Start(b)
	if err != nil {
		t.Errorf("did not expect any err: %s", err)
	}
	if p.header != 16 {
		t.Errorf("expected header: %d, got %d", 16, p.header)
	}
	if p.remainingLength != 64 {
		t.Errorf("expected remainingLength: %d, got %d", 64, p.remainingLength)
	}
	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("expected remainingLength = len(remainingBytes): %d == %d", 64, len(p.remainingBytes))
	}
	if p.PacketType() != PACKET_TYPE_CONNECT {
		t.Errorf("expected packetType %d got %d", 1, p.PacketType())
	}
	if p.Flags() != 0 {
		t.Errorf("expected flags %d got %d", 0, p.Flags())
	}
	if p.PacketLength() != 66 {
		t.Errorf("expected length %d got %d", 66, p.PacketLength())
	}
	if !p.PacketComplete() {
		t.Errorf("expected packet complete!")
	}
}
