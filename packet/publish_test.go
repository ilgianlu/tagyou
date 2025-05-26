package packet

import (
	"testing"

	"github.com/ilgianlu/tagyou/conf"
)

func TestPublish(t *testing.T) {
	p := Publish(4, 0, false, "topic", 0, []byte{0, 1, 2})

	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("Publish expected remaingLength %d, received %d", len(p.remainingBytes), p.remainingLength)
	}
	if len(p.ApplicationMessage()) != 3 {
		t.Errorf("application message expected %v, received %v", []byte{0, 1, 2}, p.ApplicationMessage())
	}
	p = Publish(4, 1, false, "topic", 123, []byte{0, 1, 2, 3, 4, 5})
	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("Publish expected remaingLength %d, received %d", len(p.remainingBytes), p.remainingLength)
	}
	if p.QoS() != 1 {
		t.Errorf("Publish expected qos 1, received %d", p.QoS())
	}
	if len(p.ApplicationMessage()) != 6 {
		t.Errorf("application message expected %v, received %v", []byte{0, 1, 2, 3, 4, 5}, p.ApplicationMessage())
	}
}

func TestParsePublish(t *testing.T) {
	// (0x30) publish packet to topic '/topic/0/messages' payload 'published 0'
	buf := []byte{0x30, 0x1e, 0x00, 0x11, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x30, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x20, 0x30}

	p, _ := Start(buf)

	if p.header != 0x30 {
		t.Errorf("expecting publish header %d, received %d", 0x30, p.header)
	}

	if p.remainingLength != len(p.remainingBytes) {
		t.Errorf("expecting remaining length %d, received %d", len(p.remainingBytes), p.remainingLength)
	}

	res := p.publishReq(conf.MQTT_V3_11)
	if res != 0 {
		t.Errorf("expecting result 0, received %d", res)
	}

	if p.PacketType() != PACKET_TYPE_PUBLISH {
		t.Errorf("expecting packet type %d, received %d", PACKET_TYPE_PUBLISH, p.PacketType())
	}

	if p.PublishTopic != "/topic/0/messages" {
		t.Errorf("expecting topic %s, received %s", "/topic/0/messages", p.PublishTopic)
	}

	if string(p.Payload()) != "published 0" {
		t.Errorf("expecting payload %s, received %s", "published 0", string(p.Payload()))
	}
}
