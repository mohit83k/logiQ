//block defines the block for the block chain and it's properties
package block

import (
	"crypto/sha256"
	"encoding/hex"
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
	s := bl.Previous + bl.Data + bl.TimeStamp + bl.Nounce
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))

}

func (bl *Block) Mine() {
	//uses set of criteria to determine if it is a valid block
	//For now : it just adds a random nounce. Extremely bad
	bl.Nounce = string(rand.Int())
	bl.Hash = bl.GetHash()

}
