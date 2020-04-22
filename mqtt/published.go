package mqtt

type Published struct {
	topic  string
	dup    bool
	qos    uint8
	retain bool
}
