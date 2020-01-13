package types

import (
	"bbc/common/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const targetBit = 16

type ProofOfWork struct {
	Block *Block //需要共识验证的区块
	//目标难度的哈希
	target *big.Int
}

//创建一个POW对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

//执行POW，比较哈希
//返回哈希值，以及碰撞次数
func (proofOfWork *ProofOfWork) Run() ([]byte, int) {
	var nonce = 0
	var hashInt big.Int
	var hash [32]byte
	//无限循环，生成符合条件的哈希值
	for {
		//生成准备数据
		dataBytes := proofOfWork.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//检测生成的哈希值是否满足条件
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			//找到了符合条件的哈希
			break
		}
		nonce++
	}
	fmt.Printf("碰撞次数：%d\n", nonce)
	return hash[:], nonce
}

//生成准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	timestampBytes := utils.Int2Hex(pow.Block.TimeStamp)
	heightBytes := utils.Int2Hex(pow.Block.Height)
	data := bytes.Join([][]byte{
		timestampBytes,
		heightBytes,
		pow.Block.PreBlockHash,
		pow.Block.HashTransaction(),
		utils.Int2Hex(nonce),
		utils.Int2Hex(targetBit),
	}, []byte{},
	)
	return data
}
