//distribute package helps to distribute blocks to all the nodes.
package distribute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mohit83k/logiQ/blockchain"
	"github.com/mohit83k/logiQ/blockchain/block"
	"github.com/mohit83k/logiQ/network/explorer"
)

const limit = 2
const blockchain_port = "37000"

func Distribute(bl block.Block) {
	jBlock, err := json.Marshal(bl)
	if err != nil {
		fmt.Println("Unable to marshal block")
		return
	}
	for _, peer := range explorer.Peers {
		url := "http://" + peer + ":" + blockchain_port + "/peer_block"

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jBlock))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		fmt.Println("Sending Block to : ", url)
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}

func RecieveBlock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RecieveBlock")
	var bl block.Block
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("No block Found"))
		return
	}

	err = json.Unmarshal(data, &bl)
	if err != nil {
		w.Write([]byte("Invalid Block Format"))
		return
	}

	if blockchain.Exists(bl) {
		w.Write([]byte("Block Added : Reached Cycle"))
		return
	}
	blockchain.AddBlock(bl)
	go Distribute(bl)
}
