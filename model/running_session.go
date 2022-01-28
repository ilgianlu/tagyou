package model

import (
	"net"
	"strings"

	"github.com/ilgianlu/tagyou/conf"
)

type RunningSession struct {
	SessionID       uint
	ClientId        string
	ProtocolVersion uint8
	LastSeen        int64
	ExpiryInterval  int64
	ConnectFlags    uint8
	KeepAlive       int
	WillTopic       string
	WillDelay       int64
	WillMessage     []byte
	Username        string
	Password        string
	SubscribeAcl    string
	PublishAcl      string
	Conn            net.Conn
}

func (s RunningSession) ReservedBit() bool {
	return (s.ConnectFlags & 0x01) == 0
}

func (s RunningSession) CleanStart() bool {
	return (s.ConnectFlags & 0x02) > 0
}

func (s RunningSession) WillFlag() bool {
	return (s.ConnectFlags & 0x04) > 0
}

func (s RunningSession) WillQoS() uint8 {
	return (s.ConnectFlags & 0x18 >> 3)
}

func (s RunningSession) WillRetain() bool {
	return (s.ConnectFlags & 0x20) > 0
}

func (s RunningSession) HavePass() bool {
	return (s.ConnectFlags & 0x40) > 0
}

func (s RunningSession) HaveUser() bool {
	return (s.ConnectFlags & 0x80) > 0
}

// func (s RunningSession) Expired() bool {
// 	return s.LastSeen+s.ExpiryInterval < time.Now().Unix()
// }

func (s RunningSession) FromLocalhost() bool {
	return strings.Index(s.Conn.RemoteAddr().String(), conf.LOCALHOST) == 0
}
