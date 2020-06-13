package mqtt

const PAYLOAD_FORMAT_INDICATOR = 0x01
const MESSAGE_EXPIRY_INTERVAL = 0x02
const CONTENT_TYPE = 0x03
const RESPONSE_TOPIC = 0x08
const CORRELATION_DATA = 0x09
const SUBSCRIPTION_IDENTIFIER = 0x0B
const SESSION_EXPIRY_INTERVAL = 0x11
const ASSIGNED_CLIENT_IDENTIFIER = 0x12
const SERVER_KEEP_ALIVE = 0x13
const AUTHENTICATION_METHOD = 0x15
const AUTHENTICATION_DATA = 0x16
const REQUEST_PROBLEM_INFORMATION = 0x17
const WILL_DELAY_INTERVAL = 0x18
const REQUEST_RESPONSE_INFORMATION = 0x19
const RESPONSE_INFORMATION = 0x1A
const SERVER_REFERENCE = 0x1C
const REASON_STRING = 0x1F
const RECEIVE_MAXIMUM = 0x21
const TOPIC_ALIAS_MAXIMUM = 0x22
const TOPIC_ALIAS = 0x23
const MAXIMUM_QOS = 0x24
const RETAIN_AVAILABLE = 0x25
const USER_PROPERTY = 0x26
const MAXIMUM_PACKET_SIZE = 0x27
const WILDCARD_SUBSCRIPTION_AVAILABLE = 0x28
const SUBSCRIPTION_IDENTIFIER_AVAILABLE = 0x29
const SHARED_SUBSCRIPTION_AVAILABLE = 0x2A

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
	for propId, prop := range p.properties {
		props = append(props, encodeProp(propId, prop)...)
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
	propertiesLength, k, err := ReadVarInt(alignProp)
	if err != nil {
		return Properties{}, 0, MALFORMED_PACKET
	}
	properties := make(map[int]Property)
	j := k
	for j < propertiesLength {
		propId, property, fw := parseProp(alignProp, j)
		if fw == 0 {
			return Properties{}, 0, MALFORMED_PACKET
		}
		properties[propId] = property
		j = j + fw
	}
	return properties, j, 0
}

func encodeProp(propId int, prop Property) []byte {
	p := WriteVarInt(propId)
	p = append(p, prop.value...)
	return p
}

func parseProp(buffer []byte, i int) (int, Property, int) {
	// tab 2.4 p.25
	propId, l, _ := ReadVarInt(buffer[i:])
	switch propId {
	case PAYLOAD_FORMAT_INDICATOR,
		REQUEST_PROBLEM_INFORMATION,
		REQUEST_RESPONSE_INFORMATION,
		MAXIMUM_QOS, RETAIN_AVAILABLE,
		WILDCARD_SUBSCRIPTION_AVAILABLE,
		SUBSCRIPTION_IDENTIFIER_AVAILABLE,
		SHARED_SUBSCRIPTION_AVAILABLE:
		// type byte
		p := Property{position: i + l}
		p.length = 1
		return propId, p, l + 1
	case MESSAGE_EXPIRY_INTERVAL,
		SESSION_EXPIRY_INTERVAL,
		WILL_DELAY_INTERVAL,
		MAXIMUM_PACKET_SIZE:
		// type 4 bytes int
		p := Property{position: i + l}
		p.length = 4
		return propId, p, l + 4
	case CONTENT_TYPE,
		RESPONSE_TOPIC,
		ASSIGNED_CLIENT_IDENTIFIER,
		AUTHENTICATION_METHOD,
		RESPONSE_INFORMATION,
		SERVER_REFERENCE,
		REASON_STRING,
		USER_PROPERTY:
		// type string
		p := Property{position: i + l + 2}
		p.length = Read2BytesInt(buffer, l)
		return propId, p, l + 2 + p.length
	case CORRELATION_DATA,
		AUTHENTICATION_DATA:
		// type binary data
		p := Property{position: i + l + 2}
		p.length = Read2BytesInt(buffer, l)
		return propId, p, l + 2 + p.length
	case SERVER_KEEP_ALIVE,
		RECEIVE_MAXIMUM,
		TOPIC_ALIAS_MAXIMUM,
		TOPIC_ALIAS:
		// type 2 bytes int
		p := Property{position: i + l}
		p.length = 2
		return propId, p, l + 2
	case SUBSCRIPTION_IDENTIFIER:
		// type var int
		p := Property{position: i + l}
		_, k, _ := ReadVarInt(buffer[l:])
		p.length = k
		return propId, p, l + k
	default:
		return 0, Property{}, 0
	}
}

func (p *Packet) getPropertyRaw(propId int) []byte {
	if prop, ok := p.properties[propId]; ok {
		return p.remainingBytes[prop.position : prop.position+prop.length]
	}
	return []byte{}
}

func (p *Packet) SessionExpiryInterval() uint32 {
	return Read4BytesInt(p.getPropertyRaw(SESSION_EXPIRY_INTERVAL))
}
