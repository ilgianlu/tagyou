package packet

type header byte

func (h header) PacketType() byte { return byte(h >> 4) }

func (h header) Flags() byte {
	return byte(h & 0x0F)
}

func (h header) QoS() byte {
	if h.PacketType() == PACKET_TYPE_PUBLISH {
		return h.Flags() & 0x06 >> 1
	}
	return 0x10
}

func (h header) Dup() bool {
	if h.PacketType() == PACKET_TYPE_PUBLISH {
		return (h.Flags() & 0x08 >> 3) == 1
	}
	return false
}

func (h header) Retain() bool {
	if h.PacketType() == PACKET_TYPE_PUBLISH {
		return (h.Flags() & 0x01) == 1
	}
	return false
}

func CheckHeader(headerByte byte) bool {
	h := header(headerByte)
	switch h.PacketType() {
	case PACKET_TYPE_CONNECT:
		if h.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_PUBLISH:
		if h.QoS() > 2 {
			return false
		}
		return true
	case PACKET_TYPE_PUBACK:
		if h.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_PUBREC:
		if h.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_PUBREL:
		if h.Flags() != 2 {
			return false
		}
		return true
	case PACKET_TYPE_PUBCOMP:
		if h.Flags() != 0 {
			return false
		}
		return true
	case PACKET_TYPE_SUBSCRIBE:
		if h.Flags() != 2 {
			return false
		}
		return true
	default:
		return true
	}
}
