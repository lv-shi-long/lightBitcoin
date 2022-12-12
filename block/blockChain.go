package block

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	genisisBlock := NewBlock("crerate by lsl", []byte{})

	blockChain := &BlockChain{Blocks: []*Block{genisisBlock}}
	return blockChain
}
