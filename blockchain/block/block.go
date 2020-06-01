//block defines the block for the block chain and it's properties
package block

type Block struct {
	index    string
	previous string
	hash     string
	data     string
	nounce   int
}
