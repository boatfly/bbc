package test

import (
	"bbc/core/types"
	"fmt"
	"testing"
)

func TestWallets_CreateWallet(t *testing.T)  {
	wallets:=types.NewWallets()
	wallets.CreateWallet()
	fmt.Printf("wallets:%v\n",wallets.Wallets)
}