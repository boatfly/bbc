package types

import (
	"bbc/common"
	"bbc/common/utils"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
)

// 区块链基本结构
type BlockChain struct {
	// Blocks []*Block //refactor to DB
	DB  *bolt.DB
	Tip []byte //保存最新区块的哈希值
}

//判断数据库文件是否存在？
func DbExist(nodeId string) bool {
	//生成不同节点的数据文件
	dbName := fmt.Sprintf(common.DbName, nodeId)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		//数据库文件不存在
		return false
	}
	return true
}

//初始化区块
func CreateBlockChainWithGenesisBlock(address string, nodeId string) *BlockChain {
	if DbExist(nodeId) {
		fmt.Println("数据库文件已存在，创世区块已存在！")
		os.Exit(1)
	}

	var lastestBlockHash []byte
	//1.创建或者打开一个数据库
	dbName := fmt.Sprintf(common.DbName, nodeId)
	db, err := bolt.Open(dbName, 0600, nil) //0600 r w x 读写执行
	if nil != err {
		log.Panicf("create db[%s] failed! %v\n", common.DbName, err)
	}

	//2.创建桶,把创世区块存入数据库中
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockTableName))
		if b == nil {
			b, err := tx.CreateBucket([]byte(common.BlockTableName))
			if nil != err {
				log.Panicf("create bucket[%s] failed! %v\n", common.BlockTableName, err)
			}

			//生成一个coinbase交易
			txCoinbase := NewCoinbaseTranscation(address) //@TODO

			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 存储
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if nil != err {
				log.Panicf("insert genesis block failed! %v\n", err)
			}
			lastestBlockHash = genesisBlock.Hash

			// 存储最新区块的哈希
			err = b.Put([]byte(common.LastestHASHKEY), genesisBlock.Hash)
			if nil != err {
				log.Panicf("insert lastesthash failed! %v\n", err)
			}
		}
		return nil
	})

	return &BlockChain{db, lastestBlockHash}
}

// 区块上链
//func (bc *BlockChain) AddBlock(txs []*Transaction) {
//	//newBlock := NewBlock(height, preBlockHash, data)
//	//bc.Blocks = append(bc.Blocks, newBlock)
//
//	//更新区块数据
//	err := bc.DB.Update(func(tx *bolt.Tx) error {
//		//获取桶
//		b := tx.Bucket([]byte(common.BlockTableName))
//		if nil != b {
//			//获取最新区块的哈希值
//			blockBytes := b.Get(bc.Tip)
//			//区块数据反序列化
//			lastestBlock := Deserialize(blockBytes)
//			//新建区块
//			newBlock := NewBlock(lastestBlock.Height+1, lastestBlock.Hash, txs)
//			//存入db
//			err := b.Put(newBlock.Hash, newBlock.Serialize())
//			if nil != err {
//				log.Panicf("insert new block failed! %v\n", err)
//			}
//			//更新最新区块哈希
//			err = b.Put([]byte(common.LastestHASHKEY), newBlock.Hash)
//			if nil != err {
//				log.Panicf("insert lastesthash failed! %v\n", err)
//			}
//
//			//更新区块链中的最新区块哈希
//			bc.Tip = newBlock.Hash
//		}
//		return nil
//	})
//	if nil != err {
//		log.Panicf("add block failed! %v\n", err)
//	}
//}

//遍历数据库，输出区块信息
func (bc *BlockChain) PrintChain() {
	fmt.Printf("读取区块完整信息")
	var currentBlock *Block
	bcit := bc.Iterator()
	for {
		fmt.Printf("........................................\n")
		currentBlock = bcit.Next()
		fmt.Printf("\tHash:%v\n", currentBlock.Hash)
		fmt.Printf("\tHeight:%d\n", currentBlock.Height)
		fmt.Printf("\tPreBlockHash:%x\n", currentBlock.PreBlockHash)
		fmt.Printf("\tTimeStamp:%v\n", currentBlock.TimeStamp)
		fmt.Printf("\tNonce:%d\n", currentBlock.Nonce)
		fmt.Printf("\tTxs:%v\n", currentBlock.Txs)
		for _, tx := range currentBlock.Txs {
			fmt.Printf("\t\ttx-hash:%v\n", tx.TxHash)

			fmt.Printf("\t\ttx-vins:\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\t\ttx-vin-TxHash:%v\n", vin.TxHash)
				fmt.Printf("\t\t\ttx-vin-Publickey:%v\n", vin.Publickey)
				fmt.Printf("\t\t\ttx-vin-Publickey-ripemd160hash:%v\n", Ripemd160Hash(vin.Publickey))
				fmt.Printf("\t\t\ttx-vin-Signature:%s\n", vin.Signature)
				fmt.Printf("\t\t\ttx-vin-Vout:%v\n", vin.Vout)
			}
			fmt.Printf("\t\ttx-vouts:\n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\t\ttx-vout-Ripemd160Hash:%v\n", vout.Ripemd160Hash)
				fmt.Printf("\t\t\ttx-vout-address_ripemd160hash:%v\n", utils.String2Ripemd160("1fmzb5HcUJqqhBQUn9tMs1c2rjwdQ3yUtB"))
				fmt.Printf("\t\t\ttx-vout-value:%v\n", vout.Value)
			}
		}
		//退出条件
		var hashInt big.Int
		hashInt.SetBytes(currentBlock.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			//发现创世区块，退出
			break
		}
	}
}

//挖矿
//通过接收交易，生成区块
func (bc *BlockChain) MineNewBlock(froms, tos, amounts []string, nodeId string) {
	var txs []*Transaction //@TODO 交易生成
	for index, _ := range froms {
		value, _ := strconv.Atoi(amounts[index])
		//生成新的交易
		tx := NewSimpleTransaction(froms[index], tos[index], value, bc, txs, nodeId)
		//@TODO 在生成交易的时候对交易进行签名
		//追加到交易列表
		txs = append(txs, tx)
		//给与交易的发起者（矿工）一定的奖励,@TODO 真是情况是矿工发起一笔交易，广播，区块上链才发起奖励，实际是在下一区块打包时考虑的操作
		rewardTx := NewCoinbaseTranscation(froms[index])
		txs = append(txs, rewardTx)
	}

	//从数据中获取最近的一个区块
	var block *Block
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockTableName))
		if nil != b {
			lastestBlockHashByte := b.Get([]byte(common.LastestHASHKEY))
			lastestBlockByte := b.Get(lastestBlockHashByte)
			block = Deserialize(lastestBlockByte)
		}
		return nil
	})
	//@TODO 在打包区块之前进行验证，对交易列表中的每一笔交易都进行验证
	for _, tx := range txs {
		//验证签名，只要有一笔失败，panic
		if bc.VerifyTransaction(tx) == false {
			log.Panicf("verify transaction failed!\n")
		}
	}

	//根据数据库中最新区块，生成新的区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	//持久化新生成的区块到数据库
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockTableName))
		if nil != b {
			err := b.Put(block.Hash, block.Serialize())
			if nil != err {
				log.Panicf("save new block to db failed!%v\n", err)
			}
			//更新区块链最新哈希值
			err = b.Put([]byte(common.LastestHASHKEY), block.Hash)
			if nil != err {
				log.Panicf("save the block lastest hash to db failed!%v\n", err)
			}
			bc.Tip = block.Hash
		}
		return nil
	})
}

//查找指定地址的可用的utxo，超过amount就中断查找
//更新当前数据库中指定地址的UTXO数量
//txs 缓存中的交易列表
func (bc *BlockChain) FindSpentableUTXOs(address string, amount int, txs []*Transaction) (int, map[string][]int) {
	spentableUTXOs := make(map[string][]int)
	var value int

	utxos := bc.FindUTXOs(address, txs)
	for _, utxo := range utxos {
		value += utxo.Output.Value
		//计算hash
		hash := hex.EncodeToString(utxo.TxHash)
		spentableUTXOs[hash] = append(spentableUTXOs[hash], utxo.Index)
		if value >= amount {
			break
		}
	}

	if amount > value {
		//余额不足
		fmt.Printf("指定钱包地址[%s]，余额[%d]不足！当前转账金额[%d]", address, value, amount)
		os.Exit(1)
	}

	return value, spentableUTXOs
}

//查找指定地址UTXO
/**
遍历查找区块链数据库中每一个区块的每一个交易
查找每一个交易中的每一个输出
判断每个输出是否满足下列条件：
1.属于传入的地址
2.是否未被花费
	1.遍历一次区块链数据库，将所有已花费的输出存入一个缓存
	2.判断当前输出输出是否在【已花费的输出缓存中】
*/
//txs 缓存中的交易列表
func (bc *BlockChain) FindUTXOs(address string, txs []*Transaction) []*UTXO {
	var uTXOs []*UTXO //当前地址的未花费输出列表
	//遍历数据库查找与from相关的交易
	//获取迭代器
	bcit := bc.Iterator()

	//获取所有已花费输出
	spentTxOutputs := bc.SpentOutputs(address)

	//缓存迭代
	//查找缓存中已花费输出
	for _, tx := range txs {
		if !tx.IsCoinbaseTransaction() {
			for _, vin := range tx.Vins {
				if vin.UnlockRipemd160Hash(utils.String2Ripemd160(address)) {
					key := hex.EncodeToString(vin.TxHash)
					//添加到已花费
					spentTxOutputs[key] = append(spentTxOutputs[key], vin.Vout)
				}
				//if vin.CheckPubkeyWithAddress(address) {
				//	key := hex.EncodeToString(vin.TxHash)
				//	//添加到已花费
				//	spentTxOutputs[key] = append(spentTxOutputs[key], vin.Vout)
				//}
			}
		}
	}
	//优先遍历缓存中的UTXO,如果余额足够，直接返回，如果余额不足，接着遍历数据库中UTXO
	//var isFull bool
	for _, tx := range txs {
	workInCache:
		for index, vout := range tx.Vouts {
			//index 当前输出在当前交易中的索引位置
			//vout 当前输出
			if vout.UnlockScriptPublicKeyWithAddress(address) {
				//if vout.CheckPubkeyWithAddress(address) {
				//判断当前输出是否已被花费
				if len(spentTxOutputs) != 0 {
					var isUtxoTx bool //判断交易是否被其他交易引用
					for txHash, indexArray := range spentTxOutputs {
						if txHash == hex.EncodeToString(tx.TxHash) {
							isUtxoTx = true
							var isSpentOutput bool
							for _, i := range indexArray {
								//txHash 当前输出所引用的交易哈希
								//indexArray txHash所关联的vout索引列表
								if index == i {
									// index==i  说明正好是当前的输出被其他交易所引用
									isSpentOutput = true
									continue workInCache
								}
							}
							if isSpentOutput == false {
								utxo := UTXO{tx.TxHash, index, vout}
								uTXOs = append(uTXOs, &utxo)
							}
						}
					}
					if isUtxoTx == false {
						utxo := UTXO{tx.TxHash, index, vout}
						uTXOs = append(uTXOs, &utxo)
					}
				} else {
					//将当前地址所有输出都添加到未花费输出列表中
					utxo := UTXO{tx.TxHash, index, vout}
					uTXOs = append(uTXOs, &utxo)
				}
			}
		}
	}

	//数据库迭代
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			//跳转
		work:
			for index, vout := range tx.Vouts {
				//index 当前输出在当前交易中的索引位置
				//vout 当前输出
				if vout.UnlockScriptPublicKeyWithAddress(address) {
					//if vout.CheckPubkeyWithAddress(address) {
					//判断当前输出是否已被花费
					if len(spentTxOutputs) != 0 {
						var isSpentOutput bool
						for txHash, indexArray := range spentTxOutputs {
							for _, i := range indexArray {
								//txHash 当前输出所引用的交易哈希
								//indexArray txHash所关联的vout索引列表
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									// index==i  说明正好是当前的输出被其他交易所引用
									isSpentOutput = true
									continue work
								}
							}
						}
						if isSpentOutput == false {
							utxo := UTXO{tx.TxHash, index, vout}
							uTXOs = append(uTXOs, &utxo)
						}
					} else {
						//将当前地址所有输出都添加到未花费输出列表中
						utxo := UTXO{tx.TxHash, index, vout}
						uTXOs = append(uTXOs, &utxo)
					}
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return uTXOs
}

//获取所有已花费输出
func (bc *BlockChain) SpentOutputs(address string) map[string][]int {
	spentTxOutputs := make(map[string][]int)
	//获取迭代器
	bcit := bc.Iterator()

	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			//排除coinbase交易
			if !tx.IsCoinbaseTransaction() {
				for _, vin := range tx.Vins {
					if vin.UnlockRipemd160Hash(utils.String2Ripemd160(address)) {
						key := hex.EncodeToString(vin.TxHash)
						//添加到已花费
						spentTxOutputs[key] = append(spentTxOutputs[key], vin.Vout)
					}
					//if vin.CheckPubkeyWithAddress(address) {
					//	key := hex.EncodeToString(vin.TxHash)
					//	//添加到已花费
					//	spentTxOutputs[key] = append(spentTxOutputs[key], vin.Vout)
					//}
				}
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return spentTxOutputs
}

//查询余额
func (bc *BlockChain) GetBalance(address string) int {
	var amount int
	utxos := bc.FindUTXOs(address, []*Transaction{})
	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}
	return amount
}

//获取一个区块兰对象
func BlockChainObject(nodeId string) *BlockChain {
	dbName := fmt.Sprintf(common.DbName, nodeId)
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("fetch db[%s] failed!%v\n", common.DbName, err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockTableName))
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

//验证签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	//coinbase交易不需要验证
	if tx.IsCoinbaseTransaction() {
		return true
	}
	preTxs := make(map[string]Transaction) //存储引用的交易
	//查找输入所引用的交易
	for _, vin := range tx.Vins {
		preTx := bc.FindTransaction(vin.TxHash)
		preTxs[hex.EncodeToString(vin.TxHash)] = preTx
	}
	return tx.Verfiy(preTxs)
}

//交易签名
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	//coinbase交易不需要签名
	if tx.IsCoinbaseTransaction() {
		return
	}
	//
	//
	preTxs := make(map[string]Transaction) //存储引用的交易
	for _, vin := range tx.Vins {
		//查找当前交易所引用的交易
		//vin.TxHash
		preTx := bc.FindTransaction(vin.TxHash)
		preTxs[hex.EncodeToString(vin.TxHash)] = preTx
	}
	tx.Sign(privateKey, preTxs)
}

//通过指定的交易哈希查找交易
func (bc *BlockChain) FindTransaction(hash []byte) Transaction {
	bcit := bc.Iterator()
	for {
		b := bcit.Next()
		for _, tx := range b.Txs {
			if bytes.Compare(tx.TxHash, hash) == 0 {
				return *tx
			}
		}
		//退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(b.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	fmt.Printf("not found tx by hash.%v\n", hash)
	return Transaction{}
}

//查找整条链上所有UTXOS
//map[string]*TxOutput key:string 交易hash
func (bc *BlockChain) FindUTXOsMap() map[string]*TxOutputs {
	//输出集合
	utxoMap := make(map[string]*TxOutputs)

	//遍历区块链
	bcit := bc.Iterator()

	//查找已花费输出
	spentTxOutputMap := bc.FindAllSpentTxOutput()

	for {
		b := bcit.Next()
		//txOutputs := &TxOutputs{[]*TxOutput{}}
		//查找所有输出
		for _, tx := range b.Txs {
			txOutputs := &TxOutputs{[]*TxOutput{}}
			txHash := hex.EncodeToString(tx.TxHash)
			//获取每笔交易的vouts
		workoutloop:
			for index, vout := range tx.Vouts {
				//查看txHash是否被【已经获取的已花费输出中】
				//获取指定交易的输入
				txInputs := spentTxOutputMap[txHash]
				if len(txInputs) > 0 {
					isSpent := false
					for _, in := range txInputs {
						//查找指定输出的所有者
						outPubkey := vout.Ripemd160Hash
						//inPubkey := in.Publickey
						//if bytes.Compare(outPubkey, Ripemd160Hash(inPubkey)) == 0 {
						//	if index == in.Vout {
						//		isSpent = true
						//		continue workoutloop
						//	}
						//}
						if in.UnlockRipemd160Hash(outPubkey) {
							if index == in.Vout {
								isSpent = true
								continue workoutloop
							}
						}
					}
					if isSpent == false {
						//当前输出没有被包含到txinputs当中
						txOutputs.TxOutputs = append(txOutputs.TxOutputs, vout)
					}
				} else {
					//没有input引用该交易，代表当前交易的所有vouts都是UTXO
					txOutputs.TxOutputs = append(txOutputs.TxOutputs, vout)
				}
			}
			utxoMap[txHash] = txOutputs
		}

		if isBreakLoop(b.PreBlockHash) {
			break
		}
	}

	return utxoMap
}

//查找整条链上所有已花费输出
func (bc *BlockChain) FindAllSpentTxOutput() map[string][]*TxInput {
	spentTxInputMap := make(map[string][]*TxInput)
	bcit := bc.Iterator()
	for {
		b := bcit.Next()
		for _, tx := range b.Txs {
			if !tx.IsCoinbaseTransaction() {
				for _, vin := range tx.Vins {
					txInputHsh := hex.EncodeToString([]byte(vin.TxHash))
					spentTxInputMap[txInputHsh] = append(spentTxInputMap[txInputHsh], vin)
				}
			}
		}
		if isBreakLoop(b.PreBlockHash) {
			break
		}
	}
	return spentTxInputMap
}

//查找整条链退出条件
func isBreakLoop(prevBlockHash []byte) bool {
	//退出循环条件
	var hashInt big.Int
	hashInt.SetBytes(prevBlockHash)
	if hashInt.Cmp(big.NewInt(0)) == 0 {
		return true
	}
	return false
}

//获取当前区块链的高度
func (bc *BlockChain) GetHeight() int64 {
	return bc.Iterator().Next().Height
}

//获取区块链当前所有的区块哈希
func (bc *BlockChain) GetBlockHashes() [][]byte {
	var blockHashes [][]byte
	for {
		b := bc.Iterator().Next()
		blockHashes = append(blockHashes, b.Hash)
		if isBreakLoop(b.PreBlockHash) {
			break
		}
	}
	return blockHashes
}

//获取指定哈希的区块
func (bc *BlockChain) GetBlock(hash []byte) []byte {
	var blockBytes []byte
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockTableName))
		if nil != b {
			blockBytes = b.Get(hash)
		}
		return nil
	})
	return blockBytes
}

func (bc *BlockChain) AddBlock(block *Block) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockTableName))
		if nil != b {
			if b.Get(block.Hash) != nil {
				//已经存在，不需要添加
				return nil
			}
			err := b.Put(block.Hash, block.Serialize())
			if nil != err {
				log.Panicf("sync block failed!%v\n", err)
			}
			blockHash := b.Get([]byte(common.LastestHASHKEY))
			lastestBlockBytes := b.Get(blockHash)
			rawBlock := Deserialize(lastestBlockBytes)
			if rawBlock.Height < block.Height {
				b.Put([]byte(common.LastestHASHKEY), block.Hash)
				bc.Tip = block.Hash
			}
		}
		return nil
	})
	if nil != err {
		log.Panicf("update db when instert the new block failed!%v\n", err)
	}
	fmt.Println("the new block is added!")
}
