package types

//交易输入管理
type TxInput struct {
	TxHash     []byte //交易哈希值
	Vout       int    //应用的上一笔的输出索引号
	ScriptSign string //用户名
}
