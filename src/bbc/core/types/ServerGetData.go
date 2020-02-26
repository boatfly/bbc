package types

//请求指定区块
type GetData struct {
	AddressFrom string //当前地址
	ID []byte //区块哈希
}
