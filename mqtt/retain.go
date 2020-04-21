package mqtt

import "time"

type Retain struct {
	topic              string
	applicationMessage []byte
	createdAt          time.Time
}
