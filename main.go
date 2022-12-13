package main

import (
	"fmt"
	"github.com/lightBitcoin/block"
)

func main() {
	fmt.Println("hello world")
	bc := block.NewBlockChain()
	bc.AddBlock("i l b c")
	bc.AddBlock("c b l i")
	for _, b := range bc.Blocks {
		b.Print()
		pow := block.NewProofOfWork(b)
		fmt.Printf("is valid:%v\n", pow.IsValid())
	}
}
