package model

type Retain struct {
	ClientID           string
	Topic              string
	ApplicationMessage []byte
	CreatedAt          int64
}
