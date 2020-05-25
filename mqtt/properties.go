type Properties map[int]int

type Property struct {
	
}

func (p *Properties) parse(p *Packet, i int) (int, int, int) {
	propertiesLength, k, err := ReadVarInt(p.remainingBytes[i:])
	if err != nil {
		return 0, 0, MALFORMED_PACKET
	}
	p.properties = make(map[int]int)
	for propertiesLength > 0 {
		propType, propLength := parseProp(p.remainingBytes[i:])
		p.properties[propType] = i
		i = i + propLength
		propertiesLength = propertiesLength - propLength
	}
	return k + propertiesLength, i, 0
}

func (p *Properties) parseProp(buffer []byte) (int, int) {
	// tab 2.4 p.25
	pType, l, _ := ReadVarInt(buffer)
	switch pType {
	case 1, 23, 25, 36, 37, 40, 41, 42:
		return pType, l + 1
	case 2, 17, 24, 39:
		return pType, l + 4
	case 3, 8, 18, 21, 26, 28, 31, 38:
		return pType, l + 2 + Read2BytesInt(buffer, l)
	case 9, 22:
		return pType, l + 2 + Read2BytesInt(buffer, l)
	case 19, 33, 34, 35:
		return pType, l + 2
	case 11:
		_, k, _ := ReadVarInt(buffer[l:])
		return pType, l + k
	default:
		return 0, 0
	}
}

func (p *Properties) get(pType int) (int, int) {
	// tab 2.4 p.25
	pType, l, _ := ReadVarInt(buffer)
	switch pType {
	case 1, 23, 25, 36, 37, 40, 41, 42:
		return pType, l + 1
	case 2, 17, 24, 39:
		return pType, l + 4
	case 3, 8, 18, 21, 26, 28, 31, 38:
		return pType, l + 2 + Read2BytesInt(buffer, l)
	case 9, 22:
		return pType, l + 2 + Read2BytesInt(buffer, l)
	case 19, 33, 34, 35:
		return pType, l + 2
	case 11:
		_, k, _ := ReadVarInt(buffer[l:])
		return pType, l + k
	default:
		return 0, 0
	}
}
