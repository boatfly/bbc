package types

import (
	"bbc/common/utils"
	"bytes"
	"fmt"
)

//交易输入管理
type TxInput struct {
	TxHash []byte //交易哈希值
	Vout   int    //应用的上一笔的输出索引号
	//ScriptSign string //用户名
	//数字签名
	Signature []byte
	Publickey []byte //公钥
}

//验证当前输入是否属于指定的钱包地址
//func (txInput *TxInput) CheckPubkeyWithAddress(address string) bool {
//	return address == txInput.ScriptSign
//}

//传递ripemd160哈希进行判断
func (txInput *TxInput) UnlockRipemd160Hash(ripemd160hash []byte) bool {
	//获取input的ripemd160的哈希值
	inputRipemd160Hash := Ripemd160Hash(txInput.Publickey)
	//2.获取校验和
	checkSumBytes := CheckSum(inputRipemd160Hash)
	//3.地址组成成员拼接
	addressBytes := append(inputRipemd160Hash, checkSumBytes...)
	//4.调用base58
	b58Bytes := utils.Base58Encode(addressBytes)
	b58Bytes=utils.String2Ripemd160(string(b58Bytes[:]))
	fmt.Printf("===1>%v\n",ripemd160hash)
	fmt.Printf("===2>%v\n",b58Bytes)
	return bytes.Compare(b58Bytes, ripemd160hash) == 0
}
