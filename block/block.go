package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/lightBitcoin/utils"
	"time"
)

// 定义 区块结构
type Block struct {
	Version uint64 // 版本协议号

	PrevBlockHash []byte // 前一个区块的HASH

	MerkleRoot []byte //梅克尔书的根

	Difficulty uint64 // 挖矿难度值

	TimeStamp uint64 // 时间撮

	Nonce uint64 //  挖矿寻找的值

	Data []byte // 包含 交易的数据

	Hash []byte // 当前区块的HASH 为了方便就写入当前区块

}

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Version:       0,
		PrevBlockHash: prevHash,
		MerkleRoot:    []byte{},
		TimeStamp:     uint64(time.Now().Unix()),
		Difficulty:    10,
		//Nonce:         10, 去掉默认的 nonce 值
		Data: []byte(data),
	}
	//block.SelfHash()
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	return block
}

// 计算自身数据的 HASH 值
func (b *Block) SelfHash() {
	//var data []byte
	//data = append(data, b.Data...)
	//data = append(data, b.PrevBlockHash...)
	//data = append(data, utils.UintToByte(b.Version)...)
	//data = append(data, b.MerkleRoot...)
	//data = append(data, utils.UintToByte(b.TimeStamp)...)
	//data = append(data, utils.UintToByte(b.Nonce)...)

	joinData := bytes.Join([][]byte{b.Data, b.PrevBlockHash,
		utils.UintToByte(b.Version),
		b.MerkleRoot,
		utils.UintToByte(b.TimeStamp),
		utils.UintToByte(b.Difficulty),
		utils.UintToByte(b.Nonce)}, []byte{})

	sum256 := sha256.Sum256(joinData)
	b.Hash = sum256[:]
}

func (b *Block) Print() {
	fmt.Println("------------------**-------------------")
	fmt.Printf("version: %d\n", b.Version)
	fmt.Printf("MerkleRoot: %x\n", b.MerkleRoot)
	fmt.Printf("Difficulty: %d\n", b.Difficulty)
	fmt.Printf("TimeStamp: %d-%s\n", b.TimeStamp, time.Unix(int64(b.TimeStamp), 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("Nonce: %d\n", b.Nonce)
	fmt.Printf("hash: %x\n", b.Hash)
	fmt.Printf("prevHash: %x\n", b.PrevBlockHash)
	fmt.Printf("data: %s\n", b.Data)
}

func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return buffer.Bytes()
}

func Deserilize(data []byte) *Block {
	//fmt.Println("接码传入的数据:%v", data)
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("decode err", err)
	}
	return &block
}
