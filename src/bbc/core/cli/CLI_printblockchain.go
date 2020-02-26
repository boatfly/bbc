package cli

import (
	"bbc/core/types"
	"fmt"
	"os"
)

//打印完整的区块信息
func (cli *CLI) printBlockChain(nodeId string) {
	//判断数据库是否存在？
	if !types.DbExist(nodeId) {
		fmt.Println("数据库不存在！")
		os.Exit(1)
	}
	//获取区块链对象实例
	blockchain := types.BlockChainObject(nodeId)
	blockchain.PrintChain()
}
