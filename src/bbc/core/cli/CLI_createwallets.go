package cli

import (
	"bbc/core/types"
	"fmt"
)

func (cli *CLI) createwallets(nodeId string)  {
	wallets:=types.NewWallets(nodeId)
	wallet:=wallets.CreateWallet(nodeId)
	fmt.Printf("wallets:%s\n",wallet.GetAddress())
}
