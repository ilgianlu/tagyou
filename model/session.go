package model

type Session struct {
	ID              int64
	LastSeen        int64
	LastConnect     int64
	ExpiryInterval  int64
	ClientId        string
	Connected       bool
	ProtocolVersion uint8
}

func (s Session) GetId() int64 {
	return s.ID
}

func (s Session) GetClientId() string {
	return s.ClientId
}

func (s Session) GetProtocolVersion() uint8 {
	return s.ProtocolVersion
}

func (s Session) Expired() bool {
	return SessionExpired(s.LastSeen, s.ExpiryInterval)
}
func (s Session) GetLastSeen() int64 {
	return s.LastSeen
}

func (s Session) GetLastConnect() int64 {
	return s.LastConnect
}
func (s Session) GetExpiryInterval() int64 {
	return s.ExpiryInterval
}
func (s Session) GetConnected() bool {
	return s.Connected
}
