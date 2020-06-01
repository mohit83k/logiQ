//explorer disovers the peers and also helps to discover the peers.
package explorer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

const limit = 2
const blockchain_port = "37000"

var Peers []string

var ip string

func init() {
	var err error
	ip, err = externalIP()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExplorerReception(w http.ResponseWriter, r *http.Request) {
	//Take the request if number of current peers is less than limit
	// send back it's IP other wise forward the request to another
	addr := strings.Split(r.RemoteAddr, ":")
	if len(Peers) < limit {
		Peers = append(Peers, ip)
		Peers = append(Peers, addr[0])
		w.Write([]byte(ip))
	}

	//Rabdomly re-direct request to other peer
	peer := Peers[rand.Intn(limit-1)]
	for peer == addr[0] {
		continue
	}
	resp, err := http.Get(peer + ":" + blockchain_port + "/explore")
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)

}

func Discover(host string) {

	if len(Peers) >= limit {
		return
	}
	resp, err := http.Get(host + ":" + blockchain_port + "/explore")
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Peers = append(Peers, string(body))
}

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
