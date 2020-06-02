package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/mohit83k/logiQ/input"
	"github.com/mohit83k/logiQ/network/distribute"
	"github.com/mohit83k/logiQ/network/explorer"
)

func main() {

	peer := flag.String("ip", "127.0.0.1", "Main Node IP address")
	flag.Parse()
	fmt.Printf("Host was %q", *peer)
	*peer = strings.TrimSpace(*peer)
	if *peer != "127.0.0.1" {
		go explorer.Discover(*peer)
	}
	fmt.Println("Needs to discover on peer : ", *peer)
	http.HandleFunc("/explore", explorer.ExplorerReception)
	http.HandleFunc("/peer_block", distribute.RecieveBlock)
	http.HandleFunc("/adddata", input.AddData)
	http.HandleFunc("/getblockchain", input.GetBlockChain)
	http.ListenAndServe(":37000", nil)

}

/**
1. Create Block and Block Chain Package
2. Create validator package
3. Crate Package for p2p network
*/
