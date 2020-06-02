package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mohit83k/logiQ/input"
	"github.com/mohit83k/logiQ/network/distribute"
	"github.com/mohit83k/logiQ/network/explorer"
)

const port = "37000"
const randomIDLength = 10
const connectURL = "/connect"
const broadcastURL = "/broadcast"
const explorerReception = "/explorer"
const distributeURL = "/distribute"
const addBlockURL = "/addblock"
const blockchainFetchURL = "/fetchblockchain"

func main() {
	rand.Seed(time.Now().UnixNano())
	peer := flag.String("peer", "", "Main Node IP address")
	//id := flag.String("id", "", "Option Feild : Identification id of the Node")
	port := flag.String("port", "37000", "Port Number for the Blockchain Node")
	local := flag.Bool("local", true, "Are all nodes running in same machine")
	flag.Parse()

	//Clean the Port Number
	if *local {
		const min = 3000
		const max = 56000
		*port = strconv.Itoa(rand.Intn(max-min) + min)
	}
	*port = strings.TrimSpace(*port)

	// if *id == "" {
	// 	*id = RandStringRunes(randomIDLength)
	// }

	selfIdentity := generateIdentity(port)

	explorer.Configure(broadcastURL, connectURL, selfIdentity)
	distribute.Configure(distributeURL)
	*peer = strings.TrimSpace(*peer)
	go func() {
		time.Sleep(5 * time.Second)
		if *peer != "" {
			explorer.ConnectNetwork(*peer)
		}
	}()

	http.HandleFunc(broadcastURL, explorer.UpdatePeerHandler)
	http.HandleFunc(connectURL, explorer.ExplorerReception)
	http.HandleFunc(addBlockURL, input.AddData)
	http.HandleFunc(blockchainFetchURL, input.GetBlockChain)
	http.HandleFunc(distributeURL, distribute.DistributeReciever)
	fmt.Println("Listening on : ", selfIdentity)
	http.ListenAndServe(":"+*port, nil)

}

// func parseUrl(url string, id *string, local *bool) string {
// 	if *local {
// 		return url
// 	}
// 	return fmt.Sprintf("/%s%s", *id, url)
// }

func generateIdentity(port *string) string {
	ip, err := externalIP()
	if err != nil {
		log.Fatal(err)
	}
	return ip + ":" + *port
}

// var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// func RandStringRunes(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 	}
// 	return string(b)
// }

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
