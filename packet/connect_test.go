package packet

import (
	"testing"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func TestConnectSuccess(t *testing.T) {
	s := model.RunningSession{}
	p := Packet{
		header:          16,
		remainingLength: 20,
		remainingBytes:  []byte{0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 0, 0, 7, 99, 108, 105, 101, 110, 116, 88},
	}
	p.connectReq(&s)
	if p.header.PacketType() != PACKET_TYPE_CONNECT {
		t.Errorf("expected packet type %d, found %d", PACKET_TYPE_CONNECT, p.header.PacketType())
	}
	if s.ProtocolVersion != conf.MQTT_V5 {
		t.Errorf("expected protocol version %d, found %d", conf.MQTT_V5, s.ProtocolVersion)
	}
}

func TestConnectSuccessWithProperties(t *testing.T) {
	s := model.RunningSession{}
	p := Packet{
		header:          16,
		remainingLength: 25,
		remainingBytes:  []byte{0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 5, 17, 0, 0, 0, 60, 0, 7, 99, 108, 105, 101, 110, 116, 88},
	}
	p.connectReq(&s)
	if p.header.PacketType() != PACKET_TYPE_CONNECT {
		t.Errorf("expected event type %d, found %d", PACKET_TYPE_CONNECT, p.header.PacketType())
	}
	if s.ProtocolVersion != conf.MQTT_V5 {
		t.Errorf("expected protocol version %d, found %d", conf.MQTT_V5, s.ProtocolVersion)
	}
	if len(p.properties) != 1 {
		t.Errorf("expected 1 property, found %d", len(p.properties))
	}
	if p.SessionExpiryInterval() != 60 {
		t.Errorf("expected sessione expiry interval 60 sec, found %d", p.SessionExpiryInterval())
	}
}
