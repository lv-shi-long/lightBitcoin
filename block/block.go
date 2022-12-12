package block

import (
	"crypto/sha256"
)

// 定义 区块结构
type Block struct {
	PrevBlockHash []byte // 前一个区块的HASH
	Hash          []byte // 当前区块的HASH 为了方便就写入当前区块
	Data          []byte // 包含 交易的数据
}

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{
		PrevBlockHash: prevHash,
		Data:          []byte(data),
	}
	block.SelfHash()
	return block
}

// 计算自身数据的 HASH 值
func (b *Block) SelfHash() {
	var data []byte
	data = append(data, b.Data...)
	data = append(data, b.PrevBlockHash...)

	sum256 := sha256.Sum256(data)
	b.Hash = sum256[:]
}
