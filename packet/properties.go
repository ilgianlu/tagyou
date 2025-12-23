package packet

import (
	"github.com/ilgianlu/tagyou/format"
)

const (
	PAYLOAD_FORMAT_INDICATOR          = 0x01
	MESSAGE_EXPIRY_INTERVAL           = 0x02
	CONTENT_TYPE                      = 0x03
	RESPONSE_TOPIC                    = 0x08
	CORRELATION_DATA                  = 0x09
	SUBSCRIPTION_IDENTIFIER           = 0x0B
	SESSION_EXPIRY_INTERVAL           = 0x11
	ASSIGNED_CLIENT_IDENTIFIER        = 0x12
	SERVER_KEEP_ALIVE                 = 0x13
	AUTHENTICATION_METHOD             = 0x15
	AUTHENTICATION_DATA               = 0x16
	REQUEST_PROBLEM_INFORMATION       = 0x17
	WILL_DELAY_INTERVAL               = 0x18
	REQUEST_RESPONSE_INFORMATION      = 0x19
	RESPONSE_INFORMATION              = 0x1A
	SERVER_REFERENCE                  = 0x1C
	REASON_STRING                     = 0x1F
	RECEIVE_MAXIMUM                   = 0x21
	TOPIC_ALIAS_MAXIMUM               = 0x22
	TOPIC_ALIAS                       = 0x23
	MAXIMUM_QOS                       = 0x24
	RETAIN_AVAILABLE                  = 0x25
	USER_PROPERTY                     = 0x26
	MAXIMUM_PACKET_SIZE               = 0x27
	WILDCARD_SUBSCRIPTION_AVAILABLE   = 0x28
	SUBSCRIPTION_IDENTIFIER_AVAILABLE = 0x29
	SHARED_SUBSCRIPTION_AVAILABLE     = 0x2A
)

type Properties map[int]Property

type Property struct {
	position int
	length   int
	value    []byte
}

func (p *Packet) encodeProperties() []byte {
	if len(p.properties) == 0 {
		return []byte{0}
	}
	props := make([]byte, 0)
	for propID, prop := range p.properties {
		props = append(props, encodeProp(propID, prop)...)
	}
	return props
}

func (p *Packet) parseProperties(i int) (int, int) {
	var j, errCode int
	p.properties, j, errCode = p.parseProps(i)
	return j, errCode
}

func (p *Packet) parseWillProperties(i int) (int, int) {
	var j, errCode int
	p.willProperties, j, errCode = p.parseProps(i)
	return j, errCode
}

func (p *Packet) parseProps(i int) (Properties, int, int) {
	alignProp := p.remainingBytes[i:]
	propertiesLength, k, err := format.ReadVarIntFromBytes(alignProp)
	if err != nil {
		return Properties{}, 0, MALFORMED_PACKET
	}
	properties := make(map[int]Property)
	j := k
	for j < propertiesLength {
		propID, property, fw := parseProp(alignProp, j, i)
		if fw == 0 {
			return Properties{}, 0, MALFORMED_PACKET
		}
		properties[propID] = property
		j = j + fw
	}
	return properties, j, 0
}

func encodeProp(propID int, prop Property) []byte {
	p, _ := format.WriteVarInt(propID)
	p = append(p, prop.value...)
	return p
}

func parseProp(buffer []byte, relativeOffset int, absoluteOffset int) (int, Property, int) {
	// tab 2.4 p.25
	propID, l, _ := format.ReadVarIntFromBytes(buffer[relativeOffset:])
	switch propID {
	case PAYLOAD_FORMAT_INDICATOR,
		REQUEST_PROBLEM_INFORMATION,
		REQUEST_RESPONSE_INFORMATION,
		MAXIMUM_QOS, RETAIN_AVAILABLE,
		WILDCARD_SUBSCRIPTION_AVAILABLE,
		SUBSCRIPTION_IDENTIFIER_AVAILABLE,
		SHARED_SUBSCRIPTION_AVAILABLE:
		// type byte
		p := Property{position: absoluteOffset + relativeOffset + l}
		p.length = 1
		return propID, p, l + 1
	case MESSAGE_EXPIRY_INTERVAL,
		SESSION_EXPIRY_INTERVAL,
		WILL_DELAY_INTERVAL,
		MAXIMUM_PACKET_SIZE:
		// type 4 bytes int
		p := Property{position: absoluteOffset + relativeOffset + l}
		p.length = 4
		return propID, p, l + 4
	case CONTENT_TYPE,
		RESPONSE_TOPIC,
		ASSIGNED_CLIENT_IDENTIFIER,
		AUTHENTICATION_METHOD,
		RESPONSE_INFORMATION,
		SERVER_REFERENCE,
		REASON_STRING,
		USER_PROPERTY:
		// type string
		p := Property{position: absoluteOffset + relativeOffset + l + 2}
		pLength, _ := format.Read2BytesInt(buffer, l)
		p.length = pLength
		return propID, p, l + 2 + p.length
	case CORRELATION_DATA, AUTHENTICATION_DATA:
		// type binary data
		p := Property{position: absoluteOffset + relativeOffset + l + 2}
		pLength, _ := format.Read2BytesInt(buffer, l)
		p.length = pLength
		return propID, p, l + 2 + p.length
	case SERVER_KEEP_ALIVE,
		RECEIVE_MAXIMUM,
		TOPIC_ALIAS_MAXIMUM,
		TOPIC_ALIAS:
		// type 2 bytes int
		p := Property{position: absoluteOffset + relativeOffset + l}
		p.length = 2
		return propID, p, l + 2
	case SUBSCRIPTION_IDENTIFIER:
		// type var int
		p := Property{position: absoluteOffset + relativeOffset + l}
		_, k, _ := format.ReadVarIntFromBytes(buffer[l:])
		p.length = k
		return propID, p, l + k
	default:
		return 0, Property{}, 0
	}
}

func (p *Packet) getPropertyRaw(propID int) []byte {
	if prop, ok := p.properties[propID]; ok {
		return p.remainingBytes[prop.position : prop.position+prop.length]
	}
	return []byte{}
}

func (p *Packet) SessionExpiryInterval() int64 {
	propRawVal := p.getPropertyRaw(SESSION_EXPIRY_INTERVAL)
	if len(propRawVal) == 0 {
		return 0
	} else {
		v, _ := format.Read4BytesInt(propRawVal)
		return int64(v)
	}
}
