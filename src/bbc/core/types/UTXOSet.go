package types

import (
	"bbc/common"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//utxo持久化相关管理
//保存指定区块链中所有的UTXO
type UTXOSet struct {
	BlockChain *BlockChain
}

//更新
func (utxoSet *UTXOSet) Refresh() {
	//获取最新区块
	lastestBlock := utxoSet.BlockChain.Iterator().Next()
	utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.UTXO_TABLE))
		if nil != b {
			//只需查找最新区块的交易列表
			for _, tx := range lastestBlock.Txs {
				if !tx.IsCoinbaseTransaction() {
					//将当前此笔交易的输入引用的utxo移除
					for _, vin := range tx.Vins {
						updateOutpus := TxOutputs{}
						//获取指定输入所引用的交易哈希的输出
						outputsByes := b.Get(vin.TxHash)
						//获取输出列表
						outs := DeserializeTXOutputs(outputsByes)
						for out_id, out := range outs.TxOutputs {
							if vin.Vout != out_id {
								updateOutpus.TxOutputs = append(updateOutpus.TxOutputs, out)
							}
						}
						if len(updateOutpus.TxOutputs) == 0 {
							b.Delete(vin.TxHash)
						} else {
							b.Put(vin.TxHash, updateOutpus.Serialize())
						}
					}

				}
				//获取当前区块中新生成的交易输出
				newOutputs := TxOutputs{}
				newOutputs.TxOutputs = append(newOutputs.TxOutputs, tx.Vouts...)
				b.Put(tx.TxHash, newOutputs.Serialize())
			}
		}
		return nil
	})
}

//查找
func (utxoSet *UTXOSet) FindUTXOwithAddress(address string) []*UTXO {
	var utxos []*UTXO
	err := utxoSet.BlockChain.DB.View(func(tx *bolt.Tx) error {
		//获取utxo table表
		b := tx.Bucket([]byte(common.UTXO_TABLE))
		if nil != b {
			//cursor 游标
			c := b.Cursor()
			//通过游标遍历boltdb数据库中的数据
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTXOutputs(v)
				for _, txoutput := range txOutputs.TxOutputs {
					if txoutput.UnlockScriptPublicKeyWithAddress(address) {
						utxo_single := &UTXO{Output: txoutput}
						utxos = append(utxos, utxo_single)
					}
				}
			}
		}
		return nil
	})
	if nil != err {
		fmt.Printf("address=%s utxo not found from utxo table \n", address)
	}
	return utxos
}

//查找余额
func (utxoSet *UTXOSet) GetBalance(address string) int {
	utxos := utxoSet.FindUTXOwithAddress(address)
	ret := 0
	if len(utxos) > 0 {
		for _, utxo := range utxos {
			fmt.Printf("utxo-txhash:%x\n", utxo.TxHash)
			fmt.Printf("utxo-index:%x\n", utxo.Index)
			fmt.Printf("utxo-Ripemd160hash:%x\n", utxo.Output.Ripemd160Hash)
			fmt.Printf("utxo-value:%d\n", utxo.Output.Value)
			ret += utxo.Output.Value
		}
	}
	return ret
}

//重置
func (utxoSet *UTXOSet) ResetUTXOSet() {
	//在第一次创建的时候更新UTXO TABLE
	utxoSet.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		//查找utxo table
		b := tx.Bucket([]byte(common.UTXO_TABLE))
		if nil != b {
			err := tx.DeleteBucket([]byte(common.UTXO_TABLE))
			if nil != err && err != bolt.ErrBucketNotFound {
				log.Panicf("reset utxo table failed!%v\n", err)
			}
		}
		bucket, err := tx.CreateBucket([]byte(common.UTXO_TABLE))
		if nil != err {
			log.Panicf("create utxo table failed!%v\n", err)
		}
		if nil != bucket {
			//查找当前所有链上的utxo
			txOutputMap := utxoSet.BlockChain.FindUTXOsMap()
			//存储
			for txhash, txoutputs := range txOutputMap {
				//将所有UTXO存入
				txHash, err := hex.DecodeString(txhash)
				if nil != err {
					log.Panic("decodestring txoutput hash failed!%v\n", err)
				}
				fmt.Printf("txhash:%s\n", txhash)
				//存入utxo table
				err = bucket.Put(txHash, txoutputs.Serialize())
				if nil != err {
					log.Panic("serialize utxo failed!%v\n", err)
				}
			}
		}
		return nil
	})
}
