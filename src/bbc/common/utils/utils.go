package utils

import (
	"bbc/common"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//实现int64转[]byte
func Int2Hex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if nil != err {
		log.Panicf("int transact to []byte failed! %v \n", err)
	}
	return buffer.Bytes()
}

//标准JSON转切片
func JSON2Slice(jsonString string) []string {
	var strSlice []string

	if err := json.Unmarshal([]byte(jsonString), &strSlice); nil != err {
		log.Panicf("json to []string failed!%v\n", err)
	}

	return strSlice
}

//string 2 ripemd160
func String2Ripemd160(address string) []byte {
	pubKeyhash := Base58Encode([]byte(address))
	hash160 := pubKeyhash[:len(pubKeyhash)-4]
	return hash160
}

//获取节点id
func GetEnvNodeId() string {
	node_id := os.Getenv(common.NODE_ID)
	if node_id == "" {
		fmt.Println("node id is not set...")
		os.Exit(1)
	}
	return node_id
}

//gob 编码
func GobEncode(data interface{}) []byte {
	var result bytes.Buffer
	enc := gob.NewEncoder(&result)
	err := enc.Encode(data)
	if nil != err {
		log.Panicf("encode the data failed!%v\n", err)
	}
	return result.Bytes()
}

//命令转换为请求([]byte)
func Command2bytes(command string) []byte {
	var bytes [common.COMMAND_LENGTH]byte
	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

//反解析，把请求中的命令解析出来
func Byte2command(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x00 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}
