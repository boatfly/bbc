package cli

import (
	"bbc/core/types"
	"fmt"
	"os"
)

//发起交易
func (cli *CLI) send(from, to, amount []string,nodeId string) {
	if !types.DbExist(nodeId) {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}

	if len(from) != len(to) {
		fmt.Println("交易参数不一致，请核查一致性...")
		os.Exit(1)
	}

	//获取区块链对象
	bc := types.BlockChainObject(nodeId)
	defer bc.DB.Close()
	//发起交易，生成新的区块
	bc.MineNewBlock(from, to, amount,nodeId)
	//刷新utxo table
	utxoSet := &types.UTXOSet{bc}
	utxoSet.Refresh()
}
