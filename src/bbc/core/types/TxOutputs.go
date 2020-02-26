package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

//存入所有输出的集合
type TxOutputs struct {
	TxOutputs []*TxOutput
}

//输出集合序列化
func (txOutputs *TxOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(txOutputs); nil != err {
		log.Panicf("serialize txouputs failed!%v\n", err)
	}
	return result.Bytes()
}

//输入集合反序列化
func DeserializeTXOutputs(txOutpusBytes []byte) *TxOutputs {
	var txOutputs TxOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutpusBytes))
	err := decoder.Decode(&txOutputs)
	if nil != err {
		log.Panicf("TxOutputs deserialize failed!%v\n", err)
	}
	return &txOutputs
}
