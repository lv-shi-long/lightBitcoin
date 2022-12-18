package block

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/lightBitcoin/transaction"
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

func NewBlockChain(miner string) *BlockChain {
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

			coinbase := transaction.NewCoinBase(miner, "i-love-china")
			genisisBlock := NewBlock([]*transaction.Transaction{coinbase}, []byte{0x00})
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

func (bc *BlockChain) AddBlock(txs []*transaction.Transaction) {
	// 上一个区块的HASH 就是 当前区块的 prevHash
	bc.db.Update(func(tx *bolt.Tx) error {
		b1 := tx.Bucket([]byte(BlockBucket))
		if b1 == nil {
			fmt.Println("bucket not exist should not be resonable, can not add block when genisis block not redady")
		} else {
			// 内存中创建一个新的区块，将其插入数据库，更新 尾区块的hash 值和 tail的值
			block := NewBlock(txs, bc.tail)
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

// 查找所有UTXO 。并且 返回所有的utxo 集合
func (bc *BlockChain) FindUtxos(addres string) []transaction.UTXOInfo {

	var utxoInfos []transaction.UTXOInfo
	//utxos := []transaction.TXOutput{}
	it := bc.NewBlockChainIterator()
	spentutxos := make(map[string][]int64)
	for {
		block := it.Next()
		if block == nil {
			fmt.Println("遍历结束")
			break
		}

		for _, tx := range block.Transactions {
			// 遍历所有的 交易输入
			for _, input := range tx.TXInputs {
				if input.Addres == addres {
					key := string(input.TXID)
					spentutxos[key] = append(spentutxos[key], input.Index)
				}
			}

			// 遍历所有的 交易输出
		OUTPUT:
			for i, output := range tx.TXOutputs {
				indexes := spentutxos[string(tx.TXID)]
				if len(indexes) != 0 {
					for _, j := range indexes {
						if int64(i) == j {
							// 说明 该笔交易的 Output中的 第i 笔 output，曾经作为某笔交易的 第j笔输入了，说明
							//  已经花出去了。
							continue OUTPUT
						}
					}
				}
				if addres == output.Addres {
					//utxos = append(utxos, output)
					utxoinfo := transaction.UTXOInfo{TXID: tx.TXID, Index: int64(i), Output: output}
					utxoInfos = append(utxoInfos, utxoinfo)
				}
			}
		}
	}
	return utxoInfos
}

func (bc *BlockChain) GetBalance(addres string) float64 {
	utxos := bc.FindUtxos(addres)
	var total float64
	for _, v := range utxos {
		total += v.Output.Value
	}
	fmt.Printf("%s 的余额为%f\n", addres, total)
	return total
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *transaction.Transaction {

	var utxos map[string][]int64
	var resValue float64

	utxos, resValue = bc.FindNeedUtxos(from, amount)
	if resValue < amount {
		fmt.Println("there is no enough money in from, new transaction failed")
		return nil
	}

	inputs := []transaction.TXInput{}
	outputs := []transaction.TXOutput{}
	// 创建一个交易，把找到的所有符合 要求的 input 包含进来。
	for txid, indexes := range utxos {
		for _, transactionIndex := range indexes {
			input := transaction.TXInput{
				TXID:   []byte(txid),
				Index:  transactionIndex,
				Addres: from,
			}
			inputs = append(inputs, input)
		}
	}

	// 创建输出，把钱转给 to 。
	output1 := transaction.TXOutput{
		Value:  amount,
		Addres: to,
	}
	outputs = append(outputs, output1)

	// 找到的输入 大于要转钱，将钱转回 付款人。找零
	if resValue-amount > 0.0 {
		output2 := transaction.TXOutput{
			Value:  resValue - amount,
			Addres: from,
		}
		outputs = append(outputs, output2)
	}

	var transaction = transaction.Transaction{
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	transaction.SetTransactionID()
	return &transaction
}

func (bc *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]int64, float64) {
	var needutxos = make(map[string][]int64)
	var resValue float64

	//  首先找到所有的 utxo 。
	utxoinfos := bc.FindUtxos(from)
	for _, utxoinfo := range utxoinfos {

	}
	it := bc.NewBlockChainIterator()
	spentutxos := make(map[string][]int64)
	for {
		block := it.Next()
		if block == nil {
			fmt.Println("遍历结束")
			break
		}

		for _, tx := range block.Transactions {
			// 遍历所有的 交易输入
			for _, input := range tx.TXInputs {
				if input.Addres == from {
					key := string(input.TXID)
					spentutxos[key] = append(spentutxos[key], input.Index)
				}
			}

			// 遍历所有的 交易输出
		OUTPUT:
			for i, output := range tx.TXOutputs {
				indexes := spentutxos[string(tx.TXID)]
				if len(indexes) != 0 {
					for _, j := range indexes {
						if int64(i) == j {
							// 说明 该笔交易的 Output中的 第i 笔 output，曾经作为某笔交易的 第j笔输入了，说明
							//  已经花出去了。
							continue OUTPUT
						}
					}
				}
				// 走到这里说明钱还没被花出去。
				if from == output.Addres {
					needutxos[string(tx.TXID)] = append(needutxos[string(tx.TXID)], int64(i))
					resValue += output.Value
					if resValue >= amount {
						fmt.Println("find  enogh money in utxo, stop find")
						return needutxos, resValue
					}
				}
			}
		}
	}
	return needutxos, resValue
}
