package packet

import (
	"errors"
)

func ReadVarInt(props []byte) (int, int, error) {
	if len(props) == 0 {
		return 0, 0, errors.New("empty slice, malformed value")
	}
	multiplier := 1
	value := 0
	i := 0
	encodedByte := props[i]
	for {
		value = value + int(encodedByte&127)*multiplier
		multiplier *= 128
		i++
		if encodedByte&128 != 0 {
			if i == 5 {
				return 0, 0, errors.New("malformed value, last byte says to continue")
			}
			encodedByte = props[i]
		} else {
			break
		}
	}
	return value, i, nil
}

func WriteVarInt(x int) []byte {
	if x > 268435455 {
		return []byte{}
	}
	divider := 128
	res := []byte{}
	for {
		encodedByte := x % divider
		x = x / divider
		if x > 0 {
			encodedByte = encodedByte | divider
		}
		res = append(res, byte(encodedByte))
		if x == 0 {
			break
		}
	}
	return res
}

func Read2BytesInt(a []byte, i int) int {
	v := int(a[i]) << 8
	i++
	return v + int(a[i])
}

func Write2BytesInt(i int) []byte {
	b := make([]byte, 2)
	b[0] = byte(i & 0xFF00 >> 8)
	b[1] = byte(i & 0x00FF)
	return b
}

func Read4BytesInt(a []byte) uint32 {
	v := uint32(a[0])
	for i := 1; i < 4; i++ {
		v = v<<8 + uint32(a[i])
	}
	return v
}
