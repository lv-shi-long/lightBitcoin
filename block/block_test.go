package block

import (
	"fmt"
	"github.com/lightBitcoin/transaction"
	"testing"
)

func (block *Block) appendTest() {
	block.Transactions = append(block.Transactions, &transaction.Transaction{
		TXID: []byte{0x01},
	})
	block.HashTransaction()
	fmt.Printf("root hash:%x\n", block.MerkleRoot)
}
func TestHashTransaction(t *testing.T) {
	block := Block{
		Transactions: []*transaction.Transaction{},
	}

	block.HashTransaction()
	fmt.Printf("root hash:%x\n", block.MerkleRoot)

	block.appendTest()

	block.appendTest()
	block.appendTest()

	block.appendTest()
	block.appendTest()
	block.appendTest()
	block.appendTest()

	block.appendTest()
	block.appendTest()
	block.appendTest()
	block.appendTest()

	block.appendTest()
	block.appendTest()
	block.appendTest()
	block.appendTest()
}
