package types

//UTXO 结构封装

type UTXO struct {
	//UTXO 对应的交易哈希
	TxHash []byte
	//在其所属交易的输出列表中的索引
	Index int
	Output *TxOutput
}