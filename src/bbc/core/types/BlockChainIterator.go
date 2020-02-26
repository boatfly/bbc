package types

import (
	"bbc/common"
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	DB          *bolt.DB //迭代目标
	CurrentHash []byte   //当前迭代目标的哈希
}

//创建一个迭代器对象
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.DB, bc.Tip}
}

func (bcit *BlockChainIterator) Next() *Block {
	var block *Block

	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(common.BlockTableName))
		if nil != b {
			currentBlockBytes := b.Get(bcit.CurrentHash)
			block = Deserialize(currentBlockBytes)
			//更新迭代器
			bcit.CurrentHash = block.PreBlockHash
		}
		return nil
	})
	if nil != err {
		log.Panicf("BlockChainIterator next occur error:%v\n", err)
	}

	return block
}
