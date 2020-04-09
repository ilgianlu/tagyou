package mqtt

import "errors"

func ReadVarInt(props []byte) (int, int, error) {
	multiplier := 1
	value := 0
	i := 0
	encodedByte := props[i]
	for ok := true; ok; ok = int(encodedByte&128) != 0 {
		value = value + int(encodedByte&127)*multiplier
		if multiplier > 128*128*128 {
			return 0, 0, errors.New("malformed value")
		}
		multiplier *= 128
		i++
		encodedByte = props[i]
	}
	return value, i, nil
}

func Read2BytesInt(a []byte, i int) int {
	v := int(a[i]) << 8
	i++
	return v + int(a[i])
}
