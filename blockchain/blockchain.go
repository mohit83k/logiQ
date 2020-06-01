package blockchain

import (
	"fmt"

	"github.com/mohit83k/logiQ/blockchain/block"
)

func say() {
	b := new(block.Block)
	b.nounce = 1
	fmt.Println(block)
}
