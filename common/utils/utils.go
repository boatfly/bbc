package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Int2Hex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if nil != err {
		log.Panicf("int transact to []byte failed! %v \n", err)
	}
	return buffer.Bytes()
}
