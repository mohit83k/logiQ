package blockchain

import (
	"fmt"
	"time"

	block "github.com/mohit83k/logiQ/blockchain/block"
)

var Blockchain []block.Block

func init() {
	//Add Genesis block
	genesis := block.Block{}
	genesis.Index = 0
	genesis.Previous = "00000000000000000000000000000000"
	genesis.Data = "Block chain for Logistics Author : Mohit Kumar" +
		"ID : github.com/mohit83k"
	genesis.Nounce = string(0)
	genesis.TimeStamp = fmt.Sprint(time.Now())
	genesis.Hash = genesis.GetHash()
	appendToChain(genesis)
}

func AddData(data string) block.Block {
	bl := getBlock(data)
	bl.Mine()
	appendToChain(bl)
	return bl
}

func getBlock(data string) block.Block {
	bl := block.Block{}
	bl.Data = data
	bl.TimeStamp = fmt.Sprint(time.Now())
	bl.Previous = Blockchain[len(Blockchain)-1].Hash
	bl.Index = Blockchain[len(Blockchain)-1].Index + 1
	return bl
}

func appendToChain(bl block.Block) {
	Blockchain = append(Blockchain, bl)
}

func AddBlock(bl block.Block) {
	if bl.Previous == Blockchain[bl.Index-1].Hash {
		appendToChain(bl)
	}
}

func Exists(bl block.Block) bool {
	if bl.Index < int64(len(Blockchain)) {
		return true
	}
	return false
}
