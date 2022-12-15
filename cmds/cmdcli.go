package cmds

import (
	"fmt"
	"github.com/lightBitcoin/block"
	"os"
)

const Usage = `    ./main addBlock "xxxx"      使用xxxx信息挖矿并且添加到区块链中
    ./main printChain           打印区块链信息
`

type CLI struct {
	BC *block.BlockChain
}

func (cli *CLI) Run() {
	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Println(Usage)
		return
	}
	switch cmds[1] {
	case "addBlock":
		if len(cmds) < 3 {
			fmt.Println("please input block info")
			os.Exit(1)
		} else {
			fmt.Println("add block to block chain")
			cli.AddBlock(cmds[2])
		}

	case "printChain":
		fmt.Println("print block chain info")
		cli.PrintChain()
	default:
		fmt.Println(Usage)
	}
}

func (cli *CLI) AddBlock(s string) {
	cli.BC.AddBlock(s)
}

func (cli *CLI) PrintChain() {
	var bc = cli.BC
	for it := bc.NewBlockChainIterator(); ; {
		blockIt := it.Next()
		if blockIt == nil {
			break
		}
		blockIt.Print()
	}
}
