package model

type Retain struct {
	Topic              string
	ApplicationMessage []byte
	CreatedAt          int64
}
