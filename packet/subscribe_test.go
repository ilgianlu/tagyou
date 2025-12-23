package packet

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func TestParseSubscribe(t *testing.T) {
	// (0x82) subscription of client 'client' to topic '/topic/#'
	buf := bytes.NewReader([]byte{0x82, 0x0d, 0x33, 0x41, 0x00, 0x08, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x23, 0x00})

	session := model.RunningSession{
		ProtocolVersion: conf.MQTT_V3_11,
		SessionID:       1,
		ClientId:        "client",
	}

	p := Packet{}
	p.Parse(bufio.NewReader(buf), &session)

	if p.header != 0x82 {
		t.Errorf("expecting subscribe header %d, received %d", 0x82, p.header)
	}

	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("expecting remaining length %d, received %d", len(p.remainingBytes), p.remainingLength)
	}

	res := p.subscribeReq(&session)
	if res != 0 {
		t.Errorf("expecting result 0, received %d", res)
	}

	if p.header.PacketType() != PACKET_TYPE_SUBSCRIBE {
		t.Errorf("expecting subscribe packet %d, received %d", PACKET_TYPE_SUBSCRIBE, p.header.PacketType())
	}

	if len(p.Subscriptions) != 1 {
		t.Errorf("expecting 1 subscription, received %d", len(p.Subscriptions))
	}

	if p.Subscriptions[0].Topic != "/topic/#" {
		t.Errorf("expecting subscription topic %s, received %s", "/topic/#", p.Subscriptions[0].Topic)
	}

	if p.Subscriptions[0].ClientId != "client" {
		t.Errorf("expecting subscription clientId %s, received %s", "client", p.Subscriptions[0].ClientId)
	}
}
