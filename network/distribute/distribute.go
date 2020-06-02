//distribute package helps to distribute blocks to all the nodes.
package distribute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mohit83k/logiQ/blockchain"
	"github.com/mohit83k/logiQ/blockchain/block"
	"github.com/mohit83k/logiQ/network/explorer"
)

var distributeURL string

func Configure(_distributeURL string) {
	distributeURL = _distributeURL
}

type DistributeResponse struct {
	Status string
	Err    error
	Index  string
}

func Distribute(bl block.Block) {
	jBlock, err := json.Marshal(bl)
	if err != nil {
		log.Println("Distribute : Unable to marshal block.Report Bug", err)
		return
	}
	for index, peer := range explorer.Peers {
		var dr = DistributeResponse{}
		if peer == "" {
			log.Println("Peer is Empty")
			continue
		}
		url := "http://" + peer + distributeURL
		//Goes to Distribute Reciever
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jBlock))
		if err != nil {
			log.Println("Distribute: Error in Post call", err)
			continue
		}
		if err = json.NewDecoder(resp.Body).Decode(&dr); err != nil {
			log.Println("Distribute: Unable to read Json from resp body : ", err)
			return
		}
		resp.Body.Close()
		fmt.Printf("%d block sent t0 %s : response is %v\n", index, peer, dr)

	}
}

func DistributeReciever(w http.ResponseWriter, r *http.Request) {
	var dr = DistributeResponse{}
	var bl block.Block
	dr.Status = "success"
	dr.Err = nil
	dr.Index = strconv.Itoa(int(bl.Index))
	// data, err := ioutil.ReadAll(r.Body)
	// w.Header().Set("Content-Type", "application/json")
	// if err != nil {
	// 	log.Println("DistributeReciever : Unable to read Body", err)
	// 	dr.Err = err
	// 	dr.Status = "Failure"
	// 	json.NewEncoder(w).Encode(dr)
	// 	return
	// }

	// err = json.Unmarshal(data, &bl)
	// if err != nil {
	// 	log.Println("DistributeReciever : Unable to read Body", err)
	// 	dr.Err = err
	// 	dr.Status = "Failure"
	// 	json.NewEncoder(w).Encode(dr)
	// 	return
	// }

	if err := json.NewDecoder(r.Body).Decode(&bl); err != nil {
		log.Println("DistributeReciever: Unable to read Json from resp body : ", err)
		return
	}

	if blockchain.Exists(bl) {
		dr.Status = "Block Already Exist"
		json.NewEncoder(w).Encode(dr)
		return
	}
	json.NewEncoder(w).Encode(dr)
	blockchain.AddBlock(bl)
}
