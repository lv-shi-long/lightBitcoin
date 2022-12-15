package block

import (
	"github.com/boltdb/bolt"
)

type BlockChainIterator struct {
	db      *bolt.DB
	current []byte
}

func (bc *BlockChain) NewBlockChainIterator() *BlockChainIterator {
	return &BlockChainIterator{
		db:      bc.db,
		current: bc.tail,
	}
}

func (it *BlockChainIterator) Next() *Block {
	var block *Block

	it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucket))
		if b == nil {
			panic("can not view when block chain is empty")
		}

		blockIndb := b.Get(it.current)
		if blockIndb == nil {
			block = nil
			return nil
		}
		block = Deserilize(blockIndb)
		it.current = block.PrevBlockHash
		return nil
	})
	return block
}
