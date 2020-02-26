package types

//节点区块展示
type Inv struct {
	AddressFrom string //当前节点的地址
	Hashes [][]byte //当前节点所拥有的的区块哈希列表
}
