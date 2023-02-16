package model

type Message struct {
	Topic       string
	Qos         byte
	Retained    bool
	Payload     string
	PayloadType byte
}
