package types

import (
	"bbc/common/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
)

//定义一个交易基本结构
type Transaction struct {
	TxHash []byte      //交易哈希
	Vins   []*TxInput  //输入列表
	Vouts  []*TxOutput //输出列表
}

//实现coinbase交易
//1.coinbase的使用：挖矿的时候；创世区块
func NewCoinbaseTranscation(address string) *Transaction {
	//coinbase的输入的特点
	// TxHash nil
	// Vout -1 为了对是否为coinbase交易进行判断
	// ScriptSign "系统奖励"
	txInput := &TxInput{[]byte(""), -1, nil, nil}

	//输出
	//value
	//ScriptPubkey or address
	//txOutput := &TxOutput{10, utils.String2Ripemd160(address)} //@TODO 收款人钱包地址
	txOutput := NewTxOutput(10, address)

	//输入、输出组装交易
	txCoinbase := &Transaction{nil, []*TxInput{txInput}, []*TxOutput{txOutput}} //:=在声明的同时进行初始化
	//生成交易哈希
	txCoinbase.TransactionHash()
	return txCoinbase
}

//生成交易哈希（交易序列化）
//refactor:不同时间生成交易的哈希值
func (tx *Transaction) TransactionHash() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(tx); nil != err {
		log.Panicf("get TransactionHash failed!%v\n", err)
	}
	//添加时间戳标识，不添加会导致所有的coinbase哈希值相同
	tm:=time.Now().UnixNano()
	//用于生成哈希原数据
	txHashBytes:=bytes.Join([][]byte{result.Bytes(),utils.Int2Hex(tm)},[]byte{})
	//生成哈希值
	txHash := sha256.Sum256(txHashBytes)

	tx.TxHash = txHash[:]
}

//生成普通转账交易
func NewSimpleTransaction(from string, to string, amonnt int, bc *BlockChain, txs []*Transaction,nodeId string) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput

	//调用可花费UTXO函数
	money, spendUTXOdic := bc.FindSpentableUTXOs(from, amonnt, txs)
	fmt.Printf("money:%d\n", money)
	//获取钱包集合兑现
	wallets := NewWallets(nodeId)
	//查找对应的钱包结构
	wallet := wallets.Wallets[from]
	for txHash, indexArray := range spendUTXOdic {
		//输入
		txHashBytes, err := hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string to []byte failed!%v\n", err)
		}
		for _, index := range indexArray {
			txInput := &TxInput{txHashBytes, index, nil, wallet.PublicKey} //@TODO
			txInputs = append(txInputs, txInput)
		}
	}

	//输出
	//txOutput := &TxOutput{amonnt, to}
	txOutput := NewTxOutput(amonnt, to)
	txOutputs = append(txOutputs, txOutput)
	//找零
	if amonnt < money { //@TODO
		//txOutput = &TxOutput{money - amonnt, from}
		txOutput := NewTxOutput(money-amonnt, from)
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足！\n")
	}

	tx := Transaction{nil, txInputs, txOutputs}
	tx.TransactionHash() //生成一笔完整的交易
	log.Printf("new transaction hash====>:%v\n",tx.TxHash)

	//对交易进行签名
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

//判断指定的交易是否coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}

//验证签名
func (tx *Transaction) Verfiy(preTxs map[string]Transaction) bool {
	//能否找到交易哈希
	for _, vin := range tx.Vins {
		if preTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("交易信息已经被篡改！")
		}
	}

	//提交相同的交易属性
	txCopy := tx.TrimedTxCopy()
	//使用相同的椭圆
	curve := elliptic.P256()
	//遍历tx输入，对每一笔输入所引用的输出进行验证
	for vin_id, vin := range tx.Vins {
		//获取关联交易
		preTx := preTxs[hex.EncodeToString(vin.TxHash)]
		//找到发送者（当前输入引用的哈希--输出的哈希）
		txCopy.Vins[vin_id].Publickey = preTx.Vouts[vin_id].Ripemd160Hash
		//由需要验证的数据生成哈希，必须要签名时候保持一致
		txCopy.TxHash = txCopy.Hash()
		//在比特币中，签名是一个数值对，r,s代表签名
		//从输入的signature中获取r.s
		//获取r,s 他们长度值相等
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		//获取公钥
		//公钥是由椭圆的x，y坐标构成
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(vin.Publickey)
		x.SetBytes(vin.Publickey[:(pubKeyLen / 2)])
		y.SetBytes(vin.Publickey[(pubKeyLen / 2):])
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		//调用验证签名核心函数
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			//log.Panicf("transcation[%x] verify failed!\n")
			return false
		}
	}
	return true
}

//交易签名
//preTxs 代表当前交易的输入所引用的所有OUTPUT所属的交易
func (tx *Transaction) Sign(privatekey ecdsa.PrivateKey, preTxs map[string]Transaction) {
	//处理输入，保证交易的正确性
	//检查tx中每一个输入所引用的j交易哈希是否包含在prevTxs中，如果没有包含，则说明交易被人修改了。
	for _, vin := range tx.Vins {
		if preTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("交易信息已经被篡改！")
		}
	}
	//提取需要签名的属性
	txCopy := tx.TrimedTxCopy()
	//处理交易副本的输入
	for vin_id, vin := range txCopy.Vins {
		//获取关联交易
		preTx := preTxs[hex.EncodeToString(vin.TxHash)]
		//找到发送者（当前输入引用的哈希--输出的哈希）
		txCopy.Vins[vin_id].Publickey = preTx.Vouts[vin_id].Ripemd160Hash
		//生成交易副本的哈希
		txCopy.TxHash = txCopy.Hash()
		//调用核心签名函数
		r, s, err := ecdsa.Sign(rand.Reader, &privatekey, txCopy.TxHash)
		if nil != err {
			log.Panicf("sign transaction failed!%v\n", err)
		}
		//组成交易签名
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[vin_id].Signature = signature
	}
}

//交易拷贝，生成一个专门用于交易签名的副本
func (tx *Transaction) TrimedTxCopy() Transaction {
	//重新组装生成一个新的交易
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	for _, vin := range tx.Vins {
		txInputs = append(txInputs, &TxInput{vin.TxHash, vin.Vout, nil, nil})
	}
	for _, vout := range tx.Vouts {
		txOutputs = append(txOutputs, &TxOutput{vout.Value, vout.Ripemd160Hash})
	}
	txCopy := Transaction{tx.TxHash, txInputs, txOutputs}
	return txCopy
}

//设置用于交易签名的哈希
func (tx *Transaction) Hash() []byte {
	txcopy := tx
	txcopy.TxHash = []byte{}
	hash := sha256.Sum256(txcopy.Serialize())
	return hash[:]
}

//交易的序列化
func (tx *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); nil != err {
		log.Panicf("serialize tx to bytes failed!%v\n", err)
	}
	return result.Bytes()
}
