package main

import (
	"fmt"
	"github.com/lightBitcoin/block"
)

func main() {
	fmt.Println("hello world")
	bc := block.NewBlockChain()
	for _, block := range bc.Blocks {
		fmt.Println(block)
	}
}
