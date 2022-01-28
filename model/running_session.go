package model

import (
	"net"
	"strings"
	"sync"

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
	Mu              sync.RWMutex
}

func (s *RunningSession) ReservedBit() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return (s.ConnectFlags & 0x01) == 0
}

func (s *RunningSession) CleanStart() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return (s.ConnectFlags & 0x02) > 0
}

func (s *RunningSession) WillFlag() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return (s.ConnectFlags & 0x04) > 0
}

func (s *RunningSession) WillQoS() uint8 {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return (s.ConnectFlags & 0x18 >> 3)
}

func (s *RunningSession) WillRetain() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return (s.ConnectFlags & 0x20) > 0
}

func (s *RunningSession) HavePass() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return (s.ConnectFlags & 0x40) > 0
}

func (s *RunningSession) HaveUser() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return (s.ConnectFlags & 0x80) > 0
}

func (s *RunningSession) FromLocalhost() bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return strings.Index(s.Conn.RemoteAddr().String(), conf.LOCALHOST) == 0
}

func (s *RunningSession) ApplyAcl(pubAcl string, subAcl string) {
	s.Mu.Lock()
	s.PublishAcl = pubAcl
	s.SubscribeAcl = subAcl
	s.Mu.Unlock()
}
