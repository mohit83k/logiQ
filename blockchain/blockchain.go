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
	genesis.TimeStamp = "Long Long Ago ...."
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
	fmt.Printf("Block Added to chain : %v\n", bl)
	Blockchain = append(Blockchain, bl)
}

//AddBlock adds block to blockchain
func AddBlock(bl block.Block) {
	fmt.Println("Add Block Previous Block : ", bl.Previous)
	fmt.Println("Add Block Previous Block : ", Blockchain[bl.Index-1].Hash)
	if bl.Previous == Blockchain[bl.Index-1].Hash {
		appendToChain(bl)
	} else {
		fmt.Println("Addblock Hash did not matched")
	}
}

//Exists func check is block exist in chain
func Exists(bl block.Block) bool {
	if bl.Index < int64(len(Blockchain)) {
		fmt.Println("Exists : Block Already exist")
		return true
	}

	fmt.Println("Exists : Block do not exists Exist ")

	return false
}
