package mqtt

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

func WriteVarInt(val int) []byte {
	if val > 268435455 {
		return []byte{}
	}
	divider := 128
	res := []byte{}
	for {
		r := val % divider
		val = val / divider
		if val == 0 {
			res = append([]byte{byte(r)}, res...)
			break
		} else {
			res = append([]byte{byte(128 + r)}, res...)
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
	b[1] = byte(i & 0xFF00 >> 8)
	b[0] = byte(i & 0x00FF)
	return b
}
