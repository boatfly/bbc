package test

import (
	"bbc/core/types"
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := types.NewWallet()
	fmt.Printf("wallet:%v\n", wallet)
	fmt.Printf("private key:%v\n", wallet.PrivateKey)
	fmt.Printf("public key:%v\n", wallet.PublicKey)
}

func TestGetAddress(t *testing.T)  {
	wallet:=types.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("address->%s\n",address)
	fmt.Printf("ori.address->%v\n",types.IsValidAddress([]byte(address)))
}