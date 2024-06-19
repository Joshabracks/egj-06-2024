package util

import (
	"bytes"
	"encoding/gob"
)

func Decode(message []byte, e any) error {
	buffer := bytes.NewBuffer(message)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(e)
	return err
}
