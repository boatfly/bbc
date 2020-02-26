package cli

import (
	"bbc/core/types"
	"fmt"
	"strconv"
)

//查询余额
func (cli *CLI) getBalance(address string,nodeId string) {
	//查找该钱包地址UTXO
	//blockchain:=types.BlockChainObject()
	//defer blockchain.DB.Close()
	//amount:=blockchain.GetBalance(address)
	//fmt.Println(address+"'s balance is:"+strconv.Itoa(amount))

	//refactor to load balance from utxo table
	blockchain := types.BlockChainObject(nodeId)
	defer blockchain.DB.Close()
	utxoSet := &types.UTXOSet{blockchain}
	amount := utxoSet.GetBalance(address)
	fmt.Println(address + "'s balance is:" + strconv.Itoa(amount))
}
