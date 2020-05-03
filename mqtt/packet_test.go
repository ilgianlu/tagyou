package mqtt

import (
	"testing"
)

func TestStartConnectComplete(t *testing.T) {
	bComplete := []byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46, 46, 0, 4, 117, 115, 101, 114, 0, 5, 112, 108, 117, 116, 111}
	p, err := Start(bComplete)
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

func TestStartConnectPartial(t *testing.T) {
	bPartial := []byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46}
	p, err := Start(bPartial)
	if err != nil {
		t.Errorf("did not expect any err: %s", err)
	}
	if p.header != 16 {
		t.Errorf("expected header: %d, got %d", 16, p.header)
	}
	if p.remainingLength != 64 {
		t.Errorf("expected remainingLength: %d, got %d", 64, p.remainingLength)
	}
	if len(p.remainingBytes) != 50 {
		t.Errorf("expected len(remainingBytes): %d got %d", 50, len(p.remainingBytes))
	}
	if p.PacketType() != PACKET_TYPE_CONNECT {
		t.Errorf("expected packetType %d got %d", 1, p.PacketType())
	}
	if p.Flags() != 0 {
		t.Errorf("expected flags %d got %d", 0, p.Flags())
	}
	if p.PacketLength() != 52 {
		t.Errorf("expected length %d got %d", 52, p.PacketLength())
	}
	if p.PacketComplete() {
		t.Errorf("expected packet complete!")
	}
}

func TestStartSubscribeComplete(t *testing.T) {
	bComplete := []byte{130, 18, 188, 226, 0, 5, 97, 47, 98, 47, 35, 0, 0, 5, 112, 105, 112, 112, 111, 0}
	p, err := Start(bComplete)
	if err != nil {
		t.Errorf("did not expect any err: %s", err)
	}
	if p.PacketType() != PACKET_TYPE_SUBSCRIBE {
		t.Errorf("expected packet type subscribe %d got %d", PACKET_TYPE_SUBSCRIBE, p.PacketType())
	}
}
