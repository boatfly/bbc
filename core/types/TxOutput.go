package types

type TxOutput struct {
	value        int    //金额
	ScriptPubkey string //用户名（UTXO所有者）
}
