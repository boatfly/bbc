package types

import (
	"bbc/common"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//请求处理

func handleVersion(request []byte, bc *BlockChain) {
	fmt.Println("the request of version handle...")
	var buff bytes.Buffer
	var data Version
	//解析请求
	dataBytes := request[common.COMMAND_LENGTH:]
	//生成version结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the version struct failed!%v\n", err)
	}
	//获取请求放的区块高度
	versionHeight := data.Height
	//获取自身节点区块高度
	height := bc.GetHeight()
	//如果当前区块高度大于versionHeight
	//将当前节点版本信息发送给请求节点
	if height > int64(versionHeight) {
		sendVersion(data.AddressFrom, bc)
	} else if height < int64(versionHeight) {
		//向发送方发起同步数据的请求
		sendGetBlocks(data.AddressFrom)
	}
}

func handleInv(request []byte, bc *BlockChain) {
	fmt.Println("the request of inv handle...")
	var buff bytes.Buffer
	var data Inv
	//解析请求
	dataBytes := request[common.COMMAND_LENGTH:]
	//生成Inv结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the inv struct failed!%v\n", err)
	}
	for _, hash := range data.Hashes {
		sendGetData(data.AddressFrom, hash)
	}
}

func handleGetData(request []byte, bc *BlockChain) {
	fmt.Println("the request of getData handle...")
	var buff bytes.Buffer
	var data GetData
	//解析请求
	dataBytes := request[common.COMMAND_LENGTH:]
	//生成GetData结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the GetData struct failed!%v\n", err)
	}
	blockBytes := bc.GetBlock(data.ID)
	sendBlock(data.AddressFrom, blockBytes)
}

func handleGetBlocks(request []byte, bc *BlockChain) {
	fmt.Println("the request of get blocks handle...")
	var buff bytes.Buffer
	var data GetBlocks
	//解析请求
	dataBytes := request[common.COMMAND_LENGTH:]
	//生成GetBlocks结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the version struct failed!%v\n", err)
	}
	//获取区块链所有的区块哈希
	blockHashes := bc.GetBlockHashes()
	sendInv(data.AddressFrom, blockHashes)
}

//接收到新区块的时候，进行处理
func handleBlock(request []byte, bc *BlockChain) {
	fmt.Println("the request of get block handle...")
	var buff bytes.Buffer
	var data BlockData
	//解析请求
	dataBytes := request[common.COMMAND_LENGTH:]
	//生成GetBlocks结构
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)

	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the GetBlock struct failed!%v\n", err)
	}

	//将接收到的区块，添加到区块链中
	blockBytes := data.Block
	block := Deserialize(blockBytes)
	bc.AddBlock(block)
	//更新utxo table
	utxoSet:=UTXOSet{bc}
	utxoSet.Refresh()
}
