package cli

import "bbc/core/types"

//初始化区块
func (cli *CLI) createBlockChain(address string,nodeId string) {
	bc := types.CreateBlockChainWithGenesisBlock(address,nodeId)
	defer bc.DB.Close()
	//设置utxo table重置操作
	utxoSet := &types.UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}
