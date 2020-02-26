package types

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	TimeStamp    int64 ``
	Hash         []byte
	PreBlockHash []byte
	Height       int64
	//Data         []byte
	//MerkleRoot []byte         //merkle树根节点哈希
	Txs   []*Transaction //交易数据（交易列表）
	Nonce int64
}

// 创建新的区块
func NewBlock(height int64, preBlockHash []byte, txs []*Transaction) *Block {
	var block Block

	block = Block{
		TimeStamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Hash:         nil,
		Txs:          txs,
		Height:       height,
	}

	//block.SetHash()

	//通过POW生成新的哈希值
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = int64(nonce)

	return &block
}

// 生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

//// 生成区块哈希
//func (b *Block) SetHash() {
//	//调用sha256实现hash生成
//	timestampBytes := utils.Int2Hex(b.TimeStamp)
//	heightBytes := utils.Int2Hex(b.Height)
//	blockBytes := bytes.Join([][]byte{
//		timestampBytes,
//		heightBytes,
//		b.PreBlockHash,
//		b.Data}, []byte{},
//	)
//	hash := sha256.Sum256(blockBytes)
//	b.Hash = hash[:]
//}

//序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	//新建编码对象
	encoder := gob.NewEncoder(&buffer)
	//序列化
	if err := encoder.Encode(block); nil != err {
		log.Panicf("serialize the block to []btye failed! %v\n", err)
	}

	return buffer.Bytes()
}

//反序列化
func Deserialize(blockBytes []byte) *Block {
	var block Block

	//新建解码对象
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&block); nil != err {
		log.Panicf("deserialize the block to []btye failed! %v\n", err)
	}
	return &block
}

//把指定区块所有交易记录序列化 (类似merkle的哈希计算方法)
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	//
	//txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	//将交易数据存入merkle树，然后生成merkle根节点
	merkletree := NewMerkleTree(txHashes)
	return merkletree.RootNode.Data
}
