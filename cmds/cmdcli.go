package cmds

import (
	"fmt"
	"github.com/lightBitcoin/block"
	"github.com/lightBitcoin/transaction"
	"os"
	"strconv"
)

const Usage = `    ./main addBlock "xxxx"       使用xxxx信息挖矿并且添加到区块链中
    ./main printChain                 打印区块链信息
	./main getBalance  xxxaddres      获取某个地址的余额
	./main send FROM TO AMOUNT MINER  转账命令，从FROM 转给 TO AMOUNT 的数量，由miner 来挖矿
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
			fmt.Println("function expired")
			//cli.AddBlock(cmds[2])
			os.Exit(1)
		}

	case "printChain":
		fmt.Println("print block chain info")
		cli.PrintChain()

	case "getBalance":
		if len(cmds) < 3 {
			fmt.Println("please input addres info")
			os.Exit(1)
		} else {
			cli.BC.GetBalance(cmds[2])
		}
	case "send":
		fmt.Println("转账命令被调用")
		if len(cmds) != 6 {
			fmt.Println("错误的输入格式。请按照 ./main FROM TO AMOUNT MINER  输入")
			os.Exit(1)
		} else {
			from := cmds[2]
			to := cmds[3]
			amout, err := strconv.ParseFloat(cmds[4], 64)
			if err != nil {
				fmt.Println("amount is not a float number")
				os.Exit(1)
			}
			miner := cmds[5]
			cli.Send(from, to, amout, miner)
		}

	default:
		fmt.Println(Usage)
	}
}

func (cli *CLI) AddBlock(txs []*transaction.Transaction) {
	cli.BC.AddBlock(txs)
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

func (cli *CLI) Send(from string, to string, amount float64, miner string) {
	// 创建挖矿交易，
	coinbase := transaction.NewCoinBase(miner, "i-love-china")

	newtx := block.NewTransaction(from, to, amount, cli.BC)

	txs := []*transaction.Transaction{}
	txs = append(txs, coinbase)
	if newtx == nil {
		fmt.Println("余额不足，交易失败")
		os.Exit(1)
	} else {
		txs = append(txs, newtx)
	}

	cli.BC.AddBlock(txs)

	fmt.Println("交易成功")
}
