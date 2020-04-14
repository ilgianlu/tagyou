package mqtt

import (
	"fmt"
	"testing"
)

func TestPublish(t *testing.T) {
	p := Publish(0, false, "topic", []byte{0, 1, 2})
	fmt.Println(p.header)
	fmt.Println(p.remainingBytes)
	rl, _, _ := ReadVarInt(p.header[1:])
	if rl != len(p.remainingBytes) {
		t.Errorf("Publish expected remaingLength %d, received %d", len(p.remainingBytes), rl)
	}
}
