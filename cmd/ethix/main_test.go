package main

import (
	"github.com/PatricioNapoli/ethix/pkg/client"
	"github.com/PatricioNapoli/ethix/pkg/indexer"
	"github.com/PatricioNapoli/ethix/pkg/store"
	"github.com/PatricioNapoli/ethix/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"
)

var addrDb store.Database[string]
var txnsDb store.Database[[]indexer.Transaction]

var idxr *indexer.Indexer

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	addrDb = store.NewGobDatabase[string]("addresses.gob")
	txnsDb = store.NewGobDatabase[[]indexer.Transaction]("transactions.gob")

	go StartMockServer()

	idxr = indexer.New("http://localhost:3333/", &addrDb, &txnsDb)
}

func TestBlock(t *testing.T) {
	block := idxr.GetCurrentBlock()
	if block != 1207 {
		t.Errorf("block is: %d", 0)
	}
}

func TestSubscribe(t *testing.T) {
	sub := idxr.Subscribe("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc5")
	if !sub {
		t.Errorf("subscription failed")
	}
}

func TestGetTxns(t *testing.T) {
	idxr = indexer.New("http://localhost:3333/", &addrDb, &txnsDb)

	idxr.Subscribe("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc5")

	txns := idxr.GetTransactions("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc5")
	if len(txns) == 0 {
		t.Errorf("transactions empty")
	}

	b, _ := utils.ToPrettyJSON(txns)
	log.Println(string(b))
}

func StartMockServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		rpcReq := []client.RPCRequest{}
		utils.FromJSON(body, &rpcReq)

		blockResponse := func() {
			data := utils.ReadFile("test/block.json")
			w.Write(data)
		}

		txnResponse := func() {
			data := utils.ReadFile("test/txns.json")

			w.Write(data)
		}

		mapReq := map[string]func(){
			"eth_blockNumber":      blockResponse,
			"eth_getBlockByNumber": txnResponse,
		}

		mapReq[rpcReq[0].RPCMethod]()
	})

	http.ListenAndServe("0.0.0.0:3333", mux)
}
