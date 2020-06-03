//explorer disovers the peers and also helps to discover the peers.
package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/mohit83k/logiQ/blockchain"
	"github.com/mohit83k/logiQ/blockchain/block"
)

var mux = &sync.Mutex{}
var Peers []string
var broadCastURL string
var connectURL string
var selfIdentity string

const protocol = "http://"

func Configure(_broadCastURL string, _connectURL string, _selfIdentity string) {
	mux.Lock()
	broadCastURL = _broadCastURL
	connectURL = _connectURL
	selfIdentity = _selfIdentity
	mux.Unlock()
}

//structure of Update Message
type UpdateRequest struct {
	ID string
}

//Response of Update Message Request
type UpdateResponse struct {
	ID     string
	Err    error
	Status string
}

//Response of Connect Request.
type ConnectResponse struct {
	Nodes      []string      //All the Avaialble Nodes on the Network. i.e Peers in our case
	Blockchain []block.Block //Existing Blockchain.
	Err        error
}

func GetServerlist(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	var servers = make([]string, len(Peers)+1)
	copy(servers, Peers)
	mux.Unlock()
	servers[len(Peers)] = selfIdentity
	fmt.Println("GetServerlist : Sending Response", servers)
	json.NewEncoder(w).Encode(servers)
}

func ExplorerReception(w http.ResponseWriter, r *http.Request) {

	var cr = ConnectResponse{}
	cr.Err = nil

	//Read Information from Peer request
	var pr = UpdateRequest{} // Peer Request have same structure as UpdateRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		cr.Err = err
		json.NewEncoder(w).Encode(cr)
		return
	}
	err = json.Unmarshal(body, &pr)
	if err != nil {
		cr.Err = err
		json.NewEncoder(w).Encode(cr)
		return
	}

	//Before adding new connection to Peer Update all existing nodes that there is new Peer
	broadCastPeers(pr)

	//Add your peers and yourself and tell it to requester
	cr.Nodes = make([]string, len(Peers)+1)
	copy(cr.Nodes, Peers)
	cr.Nodes = append(cr.Nodes, selfIdentity)
	cr.Blockchain = blockchain.Blockchain
	json.NewEncoder(w).Encode(cr)

	//Now add It in your peer list
	updatePeer(pr.ID)

}

//UpdatePeerHandler Updates takes update from it's peer when ever a new peer is added with their
func UpdatePeerHandler(w http.ResponseWriter, r *http.Request) {
	//Variable to to recive msg
	var msg = UpdateRequest{}
	var rmsg = UpdateResponse{} // Response Message
	rmsg.ID = selfIdentity
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rmsg.Err = err
		rmsg.Status = "Unable to read body of update msg"
		json.NewEncoder(w).Encode(rmsg)
		return
	}

	err = json.Unmarshal(data, &msg)
	if err != nil {
		rmsg.Err = err
		rmsg.Status = "Invalid Update Packet"
		json.NewEncoder(w).Encode(rmsg)
		return
	}
	updatePeer(msg.ID)
}

func updatePeer(peer string) {
	mux.Lock()
	Peers = append(Peers, peer)
	mux.Unlock()
}

//Broadcase Message to Peers
func broadCastPeers(msg UpdateRequest) {
	jMsg, _ := json.Marshal(msg)
	for _, peer := range Peers {
		//recieved by UpdatePeerHandler
		resp, err := http.Post(protocol+peer+broadCastURL, "application/json", bytes.NewReader(jMsg))
		if err != nil {
			log.Println(err)
			continue
		}
		//log.Println(string(ioutil.ReadAll(resp.Body)))
		resp.Body.Close()
	}
}

//Connect to Network Connects to blockchain network
func ConnectNetwork(peer string) {
	log.Println("ConnectNetwork: Trying to Connect to Network with Peer", peer)
	var connRequest = UpdateRequest{}
	connRequest.ID = selfIdentity
	jMsg, _ := json.Marshal(connRequest)
	//It goes to reception i.e ExplorerReception
	resp, err := http.Post(protocol+peer+connectURL, "application/json", bytes.NewReader(jMsg))
	if err != nil {
		log.Fatal("Unable to Make Connection with the Host. Please use different seed Node", err)
	}

	var cr = ConnectResponse{}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal("Unable to Read Body of Connection Response. Report Bug", err)
	// }
	// resp.Body.Close()

	// err = json.Unmarshal(body, &cr)
	// if err != nil {
	// 	log.Fatal("Unable to parse Body of Connection Response to Struct. Report Bug", err)
	// }

	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		log.Println("ConnectNetwork: Unable to read Json from resp body : ", err)
		return
	}

	if cr.Err != nil {
		log.Fatal("Error Returned By Peer. Report Bug", cr.Err)
	}

	//Update Blockchain
	blockchain.UpdateBlockchain(cr.Blockchain)

	//Update Peers
	for _, node := range cr.Nodes {
		updatePeer(node)
	}
	fmt.Printf("Peers Updated : %v\n", Peers)

}
