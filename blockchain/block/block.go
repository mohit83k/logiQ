//block defines the block for the block chain and it's properties
package block

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

type Block struct {
	Index     int64
	Previous  string
	Hash      string
	Data      string
	Nounce    string
	TimeStamp string
}

func (bl Block) GetHash() string {
	sum := sha256.Sum256([]byte(bl.Previous + bl.Data + bl.TimeStamp + bl.Nounce))
	return fmt.Sprint(sum)
}

func (bl *Block) Mine() {
	//uses set of criteria to determine if it is a valid block
	//For now : it just adds a random nounce. Extremely bad
	bl.Nounce = string(rand.Int())
	bl.Hash = bl.GetHash()

}
