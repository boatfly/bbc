package types

import "crypto/sha256"

//Merkle树管理
type MerkleTree struct {
	//根节点
	RootNode *MerkleNode
}

//Merkle节点结构
type MerkleNode struct {
	//左子节点
	Left *MerkleNode
	//右子节点
	Right *MerkleNode
	//数据
	Data []byte
}

//创建Merkle树
//txHashs：区块中交易哈希列表
//根节点之外的其他节点，偶数存在，当遇奇数时，最后一个复制一份组成偶数
func NewMerkleTree(txHashs [][]byte) *MerkleTree {
	var nodes []MerkleNode
	//判断交易条数
	if len(txHashs)%2 != 0 {
		txHashs = append(txHashs, txHashs[len(txHashs)-1])
	}
	//遍历所有交易数据，通过哈希生成叶子节点
	for _, data := range txHashs {
		node := MakeMerkleNode(nil, nil, data)
		nodes = append(nodes, *node)
	}
	//通过叶子节点创建父节点
	for i := 0; i < len(txHashs)/2; i++ {
		var parentNodes []MerkleNode //父节点列表
		for j := 0; j < len(nodes); j += 2 {
			node := MakeMerkleNode(&nodes[j], &nodes[j+1], nil)
			parentNodes = append(parentNodes, *node)
		}
		if len(parentNodes)%2 == 0 {
			parentNodes = append(parentNodes, parentNodes[len(parentNodes)-1])
		}
		nodes = parentNodes
	}
	mtree := MerkleTree{&nodes[0]}
	return &mtree
}

//创建Merkle节点
func MakeMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}

	//判断叶子节点
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		//非叶子节点
		prehashs := append(left.Data, right.Data...)
		hash := sha256.Sum256(prehashs)
		node.Data = hash[:]
	}
	//子节点的复制
	node.Left = left
	node.Right = right

	return node
}
