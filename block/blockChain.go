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

const (
	BlockDBName   = "blockChain.db"
	BlockBucket   = "blockBucket"
	TailBlockHash = "lastBlockHash"
)

func NewBlockChain() *BlockChain {
	db, err := bolt.Open(BlockDBName, 0666, nil)
	if err != nil {
		fmt.Println("open db err")
		return nil
	}

	var tail []byte
	db.Update(func(tx *bolt.Tx) error {
		b1 := tx.Bucket([]byte(BlockBucket))
		if b1 == nil {
			fmt.Println("bucket not exist prepare to create")
			bucket, err := tx.CreateBucket([]byte(BlockBucket))
			if err != nil {
				fmt.Println("create bucket err", err)
				return err
			}

			genisisBlock := NewBlock("this is first block", []byte{0x00})
			bucket.Put(genisisBlock.Hash, genisisBlock.Serialize())
			bucket.Put([]byte("lastBlockHash"), genisisBlock.Hash)

			//diskBlock := bucket.Get(genisisBlock.Hash)
			//blockInfo := Deserilize(diskBlock)
			//fmt.Println("解码后的数据:", blockInfo)

			tail = genisisBlock.Hash
		} else {
			tail = b1.Get([]byte(TailBlockHash))
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
	bc.db.Update(func(tx *bolt.Tx) error {
		b1 := tx.Bucket([]byte(BlockBucket))
		if b1 == nil {
			fmt.Println("bucket not exist should not be resonable, can not add block when genisis block not redady")
		} else {
			// 内存中创建一个新的区块，将其插入数据库，更新 尾区块的hash 值和 tail的值
			block := NewBlock(data, bc.tail)
			b1.Put(block.Hash, block.Serialize())
			b1.Put([]byte(TailBlockHash), block.Hash)
			bc.tail = block.Hash
		}
		return nil
	})
}

func (bc *BlockChain) Close() {
	bc.db.Close()
	bc.tail = nil
}

//func
