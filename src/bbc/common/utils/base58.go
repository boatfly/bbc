package utils

import (
	"bytes"
	"math/big"
)

//base58编码实现
//1.生成一个base58的编码基数表
//-0
//-l
//-I,O
var base58Alphet = []byte("123456789abcdefjhijkmnopqrstuvwxyzABCDEFJHJKLMNPQRSTUVWXYZ")

//编码函数
func Base58Encode(input []byte) []byte {
	var result []byte
	//将[]byte转换为big.int
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(base58Alphet)))
	//求余和商
	//添加判断，除掉的最终结果是否为0
	zero := big.NewInt(0)
	//设置余数，代表base58基数表中的位置
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		//得到的result，是一个倒叙的byte数组
		result = append(result, base58Alphet[mod.Int64()])
	}
	//反转result
	Reverse(result)

	//添加前缀1，代表公网地址
	result = append([]byte{(base58Alphet[0])}, result...)
	return result
}

//反转切片函数
func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

//解码函数
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	//去掉前缀
	zeroBytes := 1
	data := input[zeroBytes:]
	for _, b := range data {
		//查找指定字符在基数表中出现的索引
		charIndex := bytes.IndexByte(base58Alphet, b)
		//余数*58
		result.Mul(result, big.NewInt(int64(len(base58Alphet))))
		//乘积结果+mod
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	//转换为byte字节数组
	decoded := result.Bytes()
	return decoded
}