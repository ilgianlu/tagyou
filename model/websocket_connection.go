package model

import (
	"context"
	"net"

	"nhooyr.io/websocket"
)

type WebsocketConnection struct {
	Conn *websocket.Conn
}

func (c WebsocketConnection) Write(b []byte) (n int, e error) {
	err := c.Conn.Write(context.Background(), websocket.MessageBinary, b)
	return len(b), err
}

func (c WebsocketConnection) Close() error {
	return c.Conn.Close(websocket.StatusNormalClosure, "")
}

func (c WebsocketConnection) RemoteAddr() net.Addr {
	return &net.IPAddr{}
}
