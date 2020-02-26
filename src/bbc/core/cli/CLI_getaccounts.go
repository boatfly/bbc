package cli

import (
	"bbc/core/types"
	"fmt"
)

//获取地址列表
func (cli *CLI) getAccounts(nodeId string)  {
	wallets:=types.NewWallets(nodeId)
	fmt.Println("\t账号列表：")
	for key,_:=range wallets.Wallets{
		fmt.Printf("\t\t[%s]\n",key)
	}
}
