package packet

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ilgianlu/tagyou/model"
)

func TestStartConnectComplete(t *testing.T) {
	bComplete := bytes.NewReader([]byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46, 46, 0, 4, 117, 115, 101, 114, 0, 5, 112, 108, 117, 116, 111})
	session := model.RunningSession{}
	p := Packet{}
	err := p.Parse(bufio.NewReader(bComplete), &session)
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
	if p.header.PacketType() != PACKET_TYPE_CONNECT {
		t.Errorf("expected packetType %d got %d", 1, p.header.PacketType())
	}
	if p.header.Flags() != 0 {
		t.Errorf("expected flags %d got %d", 0, p.header.Flags())
	}
	if p.PacketLength() != 66 {
		t.Errorf("expected length %d got %d", 66, p.PacketLength())
	}
	if !p.PacketComplete() {
		t.Errorf("expected packet complete!")
	}
}

func TestStartConnectPartial(t *testing.T) {
	bPartial := bytes.NewReader([]byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46})
	session := model.RunningSession{}
	p := Packet{}
	err := p.Parse(bufio.NewReader(bPartial), &session)
	if err == nil {
		t.Errorf("expected erroring for few bytes")
	}
	if p.header != 16 {
		t.Errorf("expected header: %d, got %d", 16, p.header)
	}
	if p.remainingLength != 64 {
		t.Errorf("expected remainingLength: %d, got %d", 64, p.remainingLength)
	}
	if p.header.PacketType() != PACKET_TYPE_CONNECT {
		t.Errorf("expected packetType %d got %d", 1, p.header.PacketType())
	}
	if p.Flags() != 0 {
		t.Errorf("expected flags %d got %d", 0, p.Flags())
	}
	if p.PacketComplete() {
		t.Errorf("expected packet incomplete!")
	}
}

func TestStartSubscribeComplete(t *testing.T) {
	bComplete := bytes.NewReader([]byte{130, 18, 188, 226, 0, 5, 97, 47, 98, 47, 35, 0, 0, 5, 112, 105, 112, 112, 111, 0})
	session := model.RunningSession{}
	p := Packet{}
	err := p.Parse(bufio.NewReader(bComplete), &session)
	if err != nil {
		t.Errorf("did not expect any err: %s", err)
	}
	if p.header.PacketType() != PACKET_TYPE_SUBSCRIBE {
		t.Errorf("expected packet type subscribe %d got %d", PACKET_TYPE_SUBSCRIBE, p.header.PacketType())
	}
}

func TestStartVeryLongPacket(t *testing.T) {
	bComplete := bytes.NewReader([]byte{48, 170, 16, 0, 40, 69, 68, 67, 47, 100, 101, 118, 45, 97, 117, 116, 111, 47, 97, 45, 99, 108, 105, 101, 110, 116, 45, 105, 100, 47, 109, 121, 45, 97, 112, 112, 45, 110, 97, 109, 101, 47, 80, 85, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	session := model.RunningSession{}
	p := Packet{}
	err := p.Parse(bufio.NewReader(bComplete), &session)
	if err == nil {
		t.Errorf("expected erroring for few bytes")
	}
	if p.header.PacketType() != PACKET_TYPE_PUBLISH {
		t.Errorf("expected packet type publish %d got %d", PACKET_TYPE_PUBLISH, p.header.PacketType())
	}
	if p.remainingLength != 2090 {
		t.Errorf("expected packet remaining length %d got %d", 2090, p.remainingLength)
	}
}
