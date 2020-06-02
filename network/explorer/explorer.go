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
		//Peers = append(Peers, ip)
		Peers = append(Peers, addr[0])
		w.Write([]byte(ip))
		fmt.Println("Peer added : ", addr[0])
		return
	}

	//Rabdomly re-direct request to other peer
	peer := Peers[rand.Intn(limit-1)]
	for peer == addr[0] {
		fmt.Println("Same Peer")
		continue
	}
	resp, err := http.Get("http://" + peer + ":" + blockchain_port + "/explore")
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
	fmt.Println("Successfully acted as intermediate node , now peer is ", string(body))

}

func Discover(host string) {
	fmt.Println("Discover")
	if len(Peers) >= limit {
		fmt.Println("Peer limit reached")
		return
	}
	fmt.Printf("Host is %q", host)
	url := "http://" + host + ":" + blockchain_port + "/explore"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Peers = append(Peers, string(body))
	fmt.Println("Peer appended", Peers)
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
