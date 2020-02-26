package types

import (
	"bbc/common"
	"bbc/common/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

//发送请求
func sendMessage(to string, msg []byte) {
	fmt.Println("向服务器发送请求...")
	//连接上服务器
	conn, err := net.Dial(common.RPOTOCOL, to)
	if nil != err {
		log.Panicf("connection the server %x failed!%v\n", to, err)
	}
	defer conn.Close()
	//要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(msg))
	if nil != err {
		log.Panicf("add the data to conn failed!%v\n", err)
	}

}

//区块链版本验证
func sendVersion(to string,bc *BlockChain) {
	//获取当前节点的区块高度
	height := bc.GetHeight()
	//组装生成version
	versionData := Version{Height: int(height), AddressFrom: to}
	//数据的序列化
	data := utils.GobEncode(versionData)
	//将命令与版本组装成完整的请求
	request := append(utils.Command2bytes(common.CMD_VERSION), data...)
	//发送请求
	sendMessage(to, request)
}

//从指定节点同步数据
func sendGetBlocks(toaddress string) {
	//生成数据
	data := utils.GobEncode(GetBlocks{nodeAddress})
	//组装请求
	request := append(utils.Command2bytes(common.CMD_GETBLOCKS), data...)
	//发送请求
	sendMessage(toaddress, request)
}

//发送获取指定区块的请求
func sendGetData(toaddress string, hash []byte) {
	//生成数据
	data := utils.GobEncode(GetData{nodeAddress, hash})
	//组装请求
	request := append(utils.Command2bytes(common.CMD_GETDATA), data...)
	//发送请求
	sendMessage(toaddress, request)
}

//向其他节点展示
func sendInv(toaddress string, hashes [][]byte) {
	//生成数据
	data := utils.GobEncode(Inv{nodeAddress, hashes})
	//组装请求
	request := append(utils.Command2bytes(common.CMD_INV), data...)
	//发送请求
	sendMessage(toaddress, request)
}

//发送区块信息
func sendBlock(toaddress string, block []byte) {
	//生成数据
	data := utils.GobEncode(BlockData{nodeAddress, block})
	//组装请求
	request := append(utils.Command2bytes(common.CMD_BLOCK), data...)
	//发送请求
	sendMessage(toaddress, request)
}
