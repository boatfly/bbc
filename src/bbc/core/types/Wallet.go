package types

import (
	"bbc/common/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const addressCheckSumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

//创建一个钱包
func NewWallet() *Wallet {
	privatekey, publickey := newKeyPair()
	return &Wallet{privatekey, publickey}
}

//通过钱包生成公钥-私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//获取一个椭圆
	curve := elliptic.P256()
	//通过椭圆相关算法生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if nil != err {
		log.Panicf("generatekey failed!%v\n", err)
	}
	//通过私钥生成公钥
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

//生成地址
func GenerateAddress() {

}

//实现双哈希
func Ripemd160Hash(pubKey []byte) []byte {
	//1.sha256
	hash256 := sha256.New()
	hash256.Write(pubKey)
	hash := hash256.Sum(nil)
	//ripemd160
	rip160 := ripemd160.New()
	rip160.Write(hash)
	return rip160.Sum(nil)
}

//生成校验和
func CheckSum(input []byte) []byte {
	first_hash := sha256.Sum256(input)
	second_hash := sha256.Sum256(first_hash[:])
	return second_hash[:addressCheckSumLen]
}

//通过公钥获取地址
func (wallet *Wallet) GetAddress() []byte {
	//1.获取hash160
	ripemd160hash := Ripemd160Hash(wallet.PublicKey)
	//2.获取校验和
	checkSumBytes := CheckSum(ripemd160hash)
	//3.地址组成成员拼接
	addressBytes := append(ripemd160hash, checkSumBytes...)
	//4.调用base58
	b58Bytes := utils.Base58Encode(addressBytes)
	return b58Bytes
}

//判断地址有效性
func IsValidAddress(address []byte) bool {
	//地址通过base58解码(长度24)
	pubKey_checksum := utils.Base58Decode(address)
	//拆分，进行校验和
	checksumBytes := pubKey_checksum[len(pubKey_checksum)-addressCheckSumLen:]
	ripemd160hash := pubKey_checksum[:len(pubKey_checksum)-addressCheckSumLen]
	//生成
	checksum := CheckSum(ripemd160hash)
	//比较
	if bytes.Compare(checksumBytes, checksum) == 0 {
		return true
	}
	return false
}
