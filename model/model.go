package model

import (
	"bytes"
	"encoding/gob"
)

func GobEncode(e any) ([]byte, error) {
	res := bytes.Buffer{}
	enc := gob.NewEncoder(&res)
	err := enc.Encode(&e)
	if err != nil {
		return []byte{}, err
	}
	return res.Bytes(), nil
}

func GobDecode[T Auth | Retain | Retry](v []byte) (T, error) {
	valReader := bytes.NewReader(v)
	decoder := gob.NewDecoder(valReader)
	var t T
	err := decoder.Decode(&t)
	return t, err
}
