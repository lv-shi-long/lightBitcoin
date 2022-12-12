package main

import (
	"fmt"
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
		Hash:          []byte{0x00},
		Data:          []byte(data),
	}
	return block
}

func main() {
	fmt.Println("hello world")
	block := NewBlock("created by lsl", []byte{0x01, 0x02})
	fmt.Println(block)
}
