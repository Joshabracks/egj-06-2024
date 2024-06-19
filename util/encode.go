package util

import (
	"bytes"
	"encoding/gob"
)

func Encode(payload interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(payload)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
