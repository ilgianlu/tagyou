package mqtt

import "errors"

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

func Read2BytesInt(a []byte, i int) int {
	v := int(a[i]) << 8
	i++
	return v + int(a[i])
}
