package model

type Session interface {
	GetId() uint
	GetClientId() string
	GetProtocolVersion() uint8
	Expired() bool
	GetLastSeen() int64
	GetLastConnect() int64
	GetExpiryInterval() int64
	GetConnected() bool
}
