package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//定义一个交易基本结构
type Transaction struct {
	TxHash []byte      //交易哈希
	Vins   []*TxInput  //输入列表
	Vouts  []*TxOutput //输出列表
}

//实现coinbase交易
//1.coinbase的使用：挖矿的时候；创世区块
func NewCoinbaseTranscation(address string) *Transaction {
	//coinbase的输入的特点
	// TxHash nil
	// Vout -1 为了对是否为coinbase交易进行判断
	// ScriptSign "系统奖励"
	txInput := &TxInput{[]byte(""), -1, "system reward"}

	//输出
	//value
	//ScriptPubkey or address
	txOutput := &TxOutput{10, address} //@TODO 收款人钱包地址

	//输入、输出组装交易
	txCoinbase := &Transaction{nil, []*TxInput{txInput}, []*TxOutput{txOutput}} //:=在声明的同时进行初始化

	//生成交易哈希
	txCoinbase.TransactionHash()
	return txCoinbase
}

//生成交易哈希（交易序列化）
func (tx *Transaction) TransactionHash() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(tx); nil != err {
		log.Panicf("get TransactionHash failed!%v\n", err)
	}

	//生成哈希值
	txHash := sha256.Sum256(result.Bytes())

	tx.TxHash = txHash[:]
}
