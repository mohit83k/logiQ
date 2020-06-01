//block defines the block for the block chain and it's properties
package block

type Block struct {
	Index    string
	Previous string
	Hash     string
	Data     string
	Nounce   int
}
