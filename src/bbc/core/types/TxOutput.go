package types

import (
	"bbc/common/utils"
	"bytes"
)

type TxOutput struct {
	Value int //金额
	//ScriptPubkey string //用户名（UTXO所有者）
	Ripemd160Hash []byte ////用户名（UTXO所有者）
}

//验证当前输出是否属于指定的钱包地址
//func (txOutput *TxOutput) CheckPubkeyWithAddress(address string) bool {
//	return address == txOutput.ScriptPubkey
//}

//output身份验证
func (txOutput *TxOutput) UnlockScriptPublicKeyWithAddress(address string) bool {
	hash160 := utils.String2Ripemd160(address)
	return bytes.Compare(hash160, txOutput.Ripemd160Hash) == 0
}

//新建output对象
func NewTxOutput(vout int,address string) *TxOutput {
	return &TxOutput{vout,utils.String2Ripemd160(address)}
}
