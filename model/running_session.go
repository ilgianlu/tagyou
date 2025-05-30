package model

import (
	"strings"
	"sync"
	"time"

	"github.com/ilgianlu/tagyou/conf"
)

type RunningSession struct {
	ID              int64
	SessionID       int64
	ClientId        string
	ProtocolVersion uint8
	LastSeen        int64
	LastConnect     int64
	Connected       bool
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
	Router          Router
	Engine          Engine
	Conn            TagyouConn
	Mu              sync.RWMutex
}

func (s *RunningSession) ReservedBit() bool {
	return (s.ConnectFlags & 0x01) == 0
}

func (s *RunningSession) CleanStart() bool {
	return (s.ConnectFlags & 0x02) > 0
}

func (s *RunningSession) WillFlag() bool {
	return (s.ConnectFlags & 0x04) > 0
}

func (s *RunningSession) WillQoS() uint8 {
	return (s.ConnectFlags & 0x18 >> 3)
}

func (s *RunningSession) WillRetain() bool {
	return (s.ConnectFlags & 0x20) > 0
}

func (s *RunningSession) HavePass() bool {
	return (s.ConnectFlags & 0x40) > 0
}

func (s *RunningSession) HaveUser() bool {
	return (s.ConnectFlags & 0x80) > 0
}

func (s *RunningSession) FromLocalhost() bool {
	return strings.Index(s.Conn.RemoteAddr().String(), conf.LOCALHOST) == 0
}

func (s *RunningSession) GetClientId() string {
	return s.ClientId
}

func (s *RunningSession) GetProtocolVersion() uint8 {
	return s.ProtocolVersion
}

func (s *RunningSession) GetConn() TagyouConn {
	return s.Conn
}

func (s *RunningSession) GetKeepAlive() int {
	return s.KeepAlive
}

func (s *RunningSession) GetLastSeen() int64 {
	return s.LastSeen
}

func (s *RunningSession) GetLastConnect() int64 {
	return s.LastConnect
}

func (s *RunningSession) GetExpiryInterval() int64 {
	return s.ExpiryInterval
}

func (s *RunningSession) GetConnected() bool {
	return s.Connected
}

func (s *RunningSession) SetConnected(connected bool) {
	s.Connected = connected
}

func (s *RunningSession) ApplyAcl(pubAcl string, subAcl string) {
	s.PublishAcl = pubAcl
	s.SubscribeAcl = subAcl
}

func (s *RunningSession) ApplySessionId(sessionID int64) {
	s.SessionID = sessionID
}

func (s *RunningSession) Expired() bool {
	return SessionExpired(s.LastSeen, s.ExpiryInterval)
}

func SessionExpired(lastSeen int64, expiryInterval int64) bool {
	return lastSeen+expiryInterval < time.Now().Unix()
}

func (s *RunningSession) GetId() int64 {
	return s.ID
}
