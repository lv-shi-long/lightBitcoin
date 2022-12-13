package block

import (
	"fmt"
	"github.com/boltdb/bolt"
)

type BlockChain struct {
	//Blocks []*Block
	db   *bolt.DB // 存储区块的数据库句柄。
	tail []byte   // 最后一个区块的 HASH
}

func NewBlockChain() *BlockChain {
	db, err := bolt.Open("blockChain.db", 0666, nil)
	if err != nil {
		fmt.Println("open db err")
		return nil
	}

	defer db.Close()
	var tail []byte
	db.Update(func(tx *bolt.Tx) error {
		b1 := tx.Bucket([]byte("blockBucket"))
		if b1 == nil {
			fmt.Println("bucket not exist")
			bucket, err := tx.CreateBucket([]byte("blockBucket"))
			if err != nil {
				fmt.Println("create bucket err", err)
				return err
			}

			genisisBlock := NewBlock("this is first block", []byte{0x00})
			bucket.Put(genisisBlock.Hash, genisisBlock.Serialize())
			bucket.Put([]byte("lastBlockHash"), genisisBlock.Hash)

			tail = genisisBlock.Hash
		} else {
			tail = b1.Get([]byte("lastBlockHash"))

		}
		return nil
	})
	return &BlockChain{
		db:   db,
		tail: tail,
	}
}

func (bc *BlockChain) AddBlock(data string) {
	// 上一个区块的HASH 就是 当前区块的 prevHash
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	block := NewBlock(data, lastBlock.Hash)
	bc.Blocks = append(bc.Blocks, block)
}
