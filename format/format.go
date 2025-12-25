// Package format collection of basic coders/encoders
package format

import (
	"bufio"
	"bytes"
	"errors"
)

func ReadVarIntFromBytes(b []byte) (int, int, error) {
	if len(b) == 0 {
		return 0, 0, errors.New("invalid buffer length")
	}
	rd := bufio.NewReader(bytes.NewReader(b))
	return ReadVarInt(rd)
}

func ReadVarInt(props *bufio.Reader) (int, int, error) {
	multiplier := 1
	value := 0
	b, err := props.ReadByte()
	i := 1
	if err != nil {
		return 0, 0, err
	}
	for {
		value = value + int(b&127)*multiplier
		multiplier *= 128
		if b&128 != 0 {
			b, err = props.ReadByte()
			if err != nil {
				return 0, 0, err
			}
			if i == 4 {
				return 0, 0, errors.New("malformed value, last byte says to continue")
			}
			i++
		} else {
			break
		}
	}
	return value, i, nil
}

func WriteVarInt(x int) ([]byte, error) {
	if x > 268435455 {
		return []byte{}, errors.New("invalid int value")
	}
	res := []byte{}
	for {
		encodedByte := x % 128
		x = x / 128
		if x > 0 {
			encodedByte = encodedByte | 128
		}
		res = append(res, byte(encodedByte&0xFF))
		if x == 0 {
			break
		}
	}
	return res, nil
}

func Read2BytesInt(a []byte, i int) (int, error) {
	if len(a) < 2 {
		return 0, errors.New("invalid buffer length")
	}
	v := int(a[i]) << 8
	i++
	return v + int(a[i]), nil
}

func Write2BytesInt(i int) []byte {
	b := make([]byte, 2)
	b[0] = byte(i & 0xFF00 >> 8)
	b[1] = byte(i & 0x00FF)
	return b
}

func Read4BytesInt(a []byte) (uint32, error) {
	if len(a) < 4 {
		return 0, errors.New("invalid buffer length")
	}
	v := uint32(a[0])
	for i := 1; i < 4; i++ {
		v = v<<8 + uint32(a[i])
	}
	return v, nil
}
