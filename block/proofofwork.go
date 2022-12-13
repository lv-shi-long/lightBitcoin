package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	//"github.com/lightBitcoin/block"
	"github.com/lightBitcoin/utils"
	"math"
	"math/big"
)

type ProofOfWork struct {
	block *Block

	target *big.Int // 用于比较挖矿哈希值的 大数。
}

// 定义 挖矿难度，代表二进制中 前面有多少个0
const BitsDifficulty = 12

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := &ProofOfWork{
		block: block,
	}
	//  默认的挖矿难度值
	//targetStr := "0010000000000000000000000000000000000000000000000000000000000000"
	//var tmp big.Int
	//tmp.SetString(targetStr, 16)
	bigIntTmp := big.NewInt(1)
	bigIntTmp.Lsh(bigIntTmp, 256-BitsDifficulty)
	pow.target = bigIntTmp
	return pow
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	var res = [32]byte{}
	//fmt.Printf("%x\n", pow.target)
	for ; nonce < math.MaxUint64; nonce++ {
		res = sha256.Sum256(pow.prepareData(nonce))
		calcInt := big.Int{}
		calcInt.SetBytes(res[:])
		if calcInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功 hash:%x, nonce:%d\n", res, nonce)
			break
		} else {
			fmt.Printf("%x\r", res)
		}
	}
	if nonce < math.MaxUint64 {
		return res[:], nonce
	}
	fmt.Println("挖矿难度太大，寻遍所有Nonce 仍未找到合适的")
	return []byte{}, 0
}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	var b = pow.block
	//b.Nonce = nonce
	res := bytes.Join([][]byte{b.Data, b.PrevBlockHash,
		utils.UintToByte(b.Version),
		b.MerkleRoot,
		utils.UintToByte(b.TimeStamp),
		utils.UintToByte(b.Difficulty),
		utils.UintToByte(nonce)}, []byte{})
	return res
}

func (pow *ProofOfWork) IsValid() bool {
	sha := sha256.Sum256(pow.prepareData(pow.block.Nonce))

	//  将要校验的区块的 sha256转化为 Big.int 和 目标值作比较。
	tmp := big.Int{}
	tmp.SetBytes(sha[:])

	return tmp.Cmp(pow.target) == -1
}
