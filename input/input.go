package input

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mohit83k/logiQ/blockchain"
	"github.com/mohit83k/logiQ/network/distribute"
)

func AddData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	//talking data from form
	if data, ok := r.Form["data"]; ok {
		if len(data) > 0 {
			bl := blockchain.AddData(data[0])
			distribute.Distribute(bl)
		}
	}

	w.Write([]byte("Operation Complete"))
}

func GetBlockChain(w http.ResponseWriter, r *http.Request) {
	bc := blockchain.Blockchain
	json.NewEncoder(w).Encode(bc)

}
