package utils

import (
	"bytes"
	"encoding/gob"
)

func GobEncode(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func GobDecode(data []byte, rs interface{}) error {
	var readBuf bytes.Buffer
	dec := gob.NewDecoder(&readBuf)
	_, err := readBuf.Write(data)
	if err != nil {
		return err
	}
	err = dec.Decode(rs)
	if err != nil {
		return err
	}
	return nil
}