package types

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
)

//数据库名称
const dbName = "boat.bbc.db"

//表名称
const blockTableName = "bbc"

// 区块链基本结构
type BlockChain struct {
	// Blocks []*Block //refactor to DB
	DB  *bolt.DB
	Tip []byte //保存最新区块的哈希值
}

//判断数据库文件是否存在？
func DbExist() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		//数据库文件不存在
		return false
	}
	return true
}

//初始化区块
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {
	if DbExist() {
		fmt.Println("数据库文件已存在，创世区块已存在！")
		os.Exit(1)
	}

	var lastestBlockHash []byte
	//1.创建或者打开一个数据库

	db, err := bolt.Open(dbName, 0600, nil) //0600 r w x 读写执行
	if nil != err {
		log.Panicf("create db[%s] failed! %v\n", dbName, err)
	}

	//2.创建桶,把创世区块存入数据库中
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			b, err := tx.CreateBucket([]byte(blockTableName))
			if nil != err {
				log.Panicf("create bucket[%s] failed! %v\n", blockTableName, err)
			}

			//生成一个coinbase交易
			txCoinbase:=NewCoinbaseTranscation(address) //@TODO

			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 存储
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if nil != err {
				log.Panicf("insert genesis block failed! %v\n", err)
			}
			lastestBlockHash = genesisBlock.Hash

			// 存储最新区块的哈希
			err = b.Put([]byte("lastesthash"), genesisBlock.Hash)
			if nil != err {
				log.Panicf("insert lastesthash failed! %v\n", err)
			}
		}
		return nil
	})

	return &BlockChain{db, lastestBlockHash}
}

// 区块上链
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//newBlock := NewBlock(height, preBlockHash, data)
	//bc.Blocks = append(bc.Blocks, newBlock)

	//更新区块数据
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		//获取桶
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			//获取最新区块的哈希值
			blockBytes := b.Get(bc.Tip)
			//区块数据反序列化
			lastestBlock := Deserialize(blockBytes)
			//新建区块
			newBlock := NewBlock(lastestBlock.Height+1, lastestBlock.Hash, txs)
			//存入db
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if nil != err {
				log.Panicf("insert new block failed! %v\n", err)
			}
			//更新最新区块哈希
			err = b.Put([]byte("lastesthash"), newBlock.Hash)
			if nil != err {
				log.Panicf("insert lastesthash failed! %v\n", err)
			}

			//更新区块链中的最新区块哈希
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if nil != err {
		log.Panicf("add block failed! %v\n", err)
	}
}

//遍历数据库，输出区块信息
func (bc *BlockChain) PrintChain() {
	fmt.Printf("读取区块完整信息")
	var currentBlock *Block
	bcit := bc.Iterator()
	for {
		fmt.Printf("........................................\n")
		currentBlock = bcit.Next()
		fmt.Printf("\tHash:%x\n", currentBlock.Hash)
		fmt.Printf("\tHeight:%d\n", currentBlock.Height)
		fmt.Printf("\tPreBlockHash:%x\n", currentBlock.PreBlockHash)
		fmt.Printf("\tTxs:%v\n", currentBlock.Txs)
		fmt.Printf("\tTimeStamp:%v\n", currentBlock.TimeStamp)
		fmt.Printf("\tNonce:%d\n", currentBlock.Nonce)
		//退出条件
		var hashInt big.Int
		hashInt.SetBytes(currentBlock.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			//发现创世区块，退出
			break
		}
	}
}

//获取一个区块兰对象
func BlockChainObject() *BlockChain {
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("fetch db[%s] failed!%v\n", dbName, err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			tip = b.Get([]byte("lastesthash"))
		}
		return nil
	})
	if nil != err {
		log.Panicf("fetch blockchain failed!%v\n", err)
	}
	return &BlockChain{DB: db, Tip: tip}
}
