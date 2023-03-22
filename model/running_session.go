package model

import (
	"strings"
	"sync"
	"time"

	"github.com/ilgianlu/tagyou/conf"
)

type RunningSession struct {
	ID              uint
	SessionID       uint
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
	Conn            TagyouConn
	Mu              sync.RWMutex
}

func (s *RunningSession) ReservedBit() bool {
	return (s.ConnectFlags & 0x01) == 0
}

func (s *RunningSession) CleanStart() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
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
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.ClientId
}

func (s *RunningSession) GetProtocolVersion() uint8 {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.ProtocolVersion
}

func (s *RunningSession) GetConn() TagyouConn {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.Conn
}

func (s *RunningSession) GetKeepAlive() int {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.KeepAlive
}

func (s *RunningSession) ApplyAcl(pubAcl string, subAcl string) {
	s.Mu.Lock()
	s.PublishAcl = pubAcl
	s.SubscribeAcl = subAcl
	s.Mu.Unlock()
}

func (s *RunningSession) ApplySessionId(sessionID uint) {
	s.Mu.Lock()
	s.SessionID = sessionID
	s.Mu.Unlock()
}

func (s *RunningSession) Expired() bool {
	return SessionExpired(s.LastSeen, s.ExpiryInterval)
}

func SessionExpired(lastSeen int64, expiryInterval int64) bool {
	return lastSeen+expiryInterval < time.Now().Unix()
}
