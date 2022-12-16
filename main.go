package main

import (
	"github.com/lightBitcoin/block"
	"github.com/lightBitcoin/cmds"
)

func main() {

	//fmt.Println("hello world")

	bc := block.NewBlockChain("abcdefg")
	cli := cmds.CLI{BC: bc}
	cli.Run()
	////bc.AddBlock("1234")
	////bc.AddBlock("9876")
	defer bc.Close()
	//

}
