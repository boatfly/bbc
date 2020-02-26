package types

import (
	"bbc/common"
	"bbc/common/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//网络服务文件管理

//3001作为引导节点（主节点）
var knownNodes = []string{"localhost:3001"}
var nodeAddress string

//启动服务
func StartServer(nodeId string) {
	fmt.Printf("启动节点[%s]...\n", nodeId)
	nodeAddress = fmt.Sprintf("localhost:%s", nodeId)
	//监听节点
	listen, err := net.Listen(common.RPOTOCOL, nodeAddress)
	if nil != err {
		log.Panicf("listen address of %s failed!%v\n", nodeAddress, err)
	}
	defer listen.Close()

	//获取区块链
	bc := BlockChainObject(nodeId)

	//两个节点，主节点负责保存数据，钱包节点负责发送请求，同步数据
	if nodeAddress != knownNodes[0] {
		//不是主节点，发送请求，同步数据
		//sendMessage(knownNodes[0], nodeAddress)
		sendVersion(knownNodes[0],bc)
	}

	for {
		//生成连接，接收请求
		conn, err := listen.Accept()
		if nil != err {
			log.Panicf("accept connect failed!%v\n", err)
		}
		//处理请求
		//单独启动一个协程（goroutine），进行请求处理
		handlerConnection(conn, bc)
	}
}

//请求处理函数
func handlerConnection(conn net.Conn, bc *BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if nil != err {
		log.Panicf("receive message failed!%v\n", err)
	}
	fmt.Printf("receive a message %v\n", request)
	cmd := utils.Byte2command(request[0:common.COMMAND_LENGTH])
	fmt.Printf("receive CMD IS: %s\n", cmd)
	switch cmd {
	case common.CMD_VERSION:
		handleVersion(request, bc)
	case common.CMD_GETBLOCKS:
		handleGetBlocks(request, bc)
	case common.CMD_BLOCK:
		handleBlock(request, bc)
	case common.CMD_GETDATA:
		handleGetData(request, bc)
	case common.CMD_INV:
		handleInv(request, bc)
	default:
		fmt.Printf("unknow command!")
	}
}
