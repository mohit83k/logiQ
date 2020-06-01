package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mohit83k/logiQ/blockchain"
	"github.com/mohit83k/logiQ/network/distribute"
	"github.com/mohit83k/logiQ/network/explorer"
)

func main() {
	//Initiating web calls
	go func() {
		http.HandleFunc("/explore", explorer.ExplorerReception)
		http.HandleFunc("/peer_block", distribute.RecieveBlock)
		http.ListenAndServe(":37000", nil)
	}()

	peer := flag.String("ip", "127.0.0.1", "Main Node IP address")
	flag.Parse()
	if *peer != "127.0.0.1" {
		go explorer.Discover(*peer)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Transaction Record: ")
		text, _ := reader.ReadString('\n')
		if strings.TrimRight(text, "\n") == "print_blockchain" {
			fmt.Printf("%v\n", blockchain.Blockchain)
		} else {
			bl := blockchain.AddData(text)
			fmt.Printf("%v\n", bl)
			go distribute.Distribute(bl)
		}

	}
}

/**
1. Create Block and Block Chain Package
2. Create validator package
3. Crate Package for p2p network
*/
