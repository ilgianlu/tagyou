package model

import "net"

type TagyouConn interface {
	Write(b []byte) (n int, err error)
	Close() error
	RemoteAddr() net.Addr
}
