package model

import (
	"time"
)

type Session struct {
	ID              uint
	LastSeen        int64
	LastConnect     int64
	ExpiryInterval  int64
	ClientId        string
	Connected       bool
	ProtocolVersion uint8
	Subscriptions   []Subscription `json:"-"`
	Retries         []Retry        `json:"-"`
}

func (s Session) Expired() bool {
	return s.LastSeen+s.ExpiryInterval < time.Now().Unix()
}

func (s *Session) UpdateFromRunning(running *RunningSession) {
	running.Mu.RLock()
	s.ProtocolVersion = running.ProtocolVersion
	s.ExpiryInterval = running.ExpiryInterval
	s.LastConnect = running.LastConnect
	s.Connected = true
	running.Mu.RUnlock()
}
