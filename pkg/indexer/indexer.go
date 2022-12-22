package indexer

import (
	"fmt"
	"github.com/PatricioNapoli/ethix/pkg/client"
	"github.com/PatricioNapoli/ethix/pkg/store"
	"github.com/PatricioNapoli/ethix/pkg/utils"
	"log"
	"math/big"
	"time"
)

type Parser interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock() int

	// Subscribe add address to observer
	Subscribe(address string) bool

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}

type Indexer struct {
	RpcEndpoint    string
	AddressesDb    *store.Database[string]
	TransactionsDb *store.Database[[]Transaction]
	LastBlock      int
}

func New(endpoint string, addressesDb *store.Database[string], transactionsDb *store.Database[[]Transaction]) *Indexer {
	// TODO: create coroutine that checks if current block changed and if it did, send an RPC request of current block
	//       number transactions, then find matches out or in with addresses db, and add them to transactions db.

	idxr := &Indexer{
		RpcEndpoint:    endpoint,
		AddressesDb:    addressesDb,
		TransactionsDb: transactionsDb,
		LastBlock:      0,
	}

	idxr.index()
	go utils.DoEvery(5*time.Second, idxr.index)

	return idxr
}

func (i *Indexer) index() {
	curr := i.GetCurrentBlock()
	if curr <= i.LastBlock {
		return
	}

	i.LastBlock = curr

	err, res := client.NewRpc[Block](i.RpcEndpoint, "2.0").PostRequests([]*client.RPCRequest{
		{
			RPCMethod: "eth_getBlockByNumber",
			Params:    []interface{}{fmt.Sprintf("0x%x", i.LastBlock), true},
		},
	})

	if err != nil {
		log.Printf("error when indexing: %s", err.Error())
		return
	}

	if len(res) == 0 {
		log.Printf("empty indexer request on block number %d: %s", i.LastBlock, err.Error())
		return
	}

	tempIndex := map[string]Transaction{}

	r := res[0].Result
	txns := r.Transactions

	for _, txn := range txns {
		tempIndex[txn.From] = txn
		tempIndex[txn.To] = txn
	}

	addressesTable := (*i.AddressesDb).ReadOrCreate("eth_addresses")
	addressesTable.Lock()
	defer addressesTable.Unlock()

	for k, _ := range addressesTable.Data {
		if txn, exists := tempIndex[k]; exists {
			txnTable := (*i.TransactionsDb).ReadOrCreate("eth_transactions")
			txnList := txnTable.Get(k)
			txnList = append(txnList, txn)
			txnTable.Set(k, txnList)
		}
	}
}

func (i *Indexer) GetCurrentBlock() int {
	err, res := client.NewRpc[string](i.RpcEndpoint, "2.0").PostRequests([]*client.RPCRequest{
		{
			RPCMethod: "eth_blockNumber",
			Params:    []interface{}{},
		},
	})

	if err != nil {
		panic(err.Error())
	}

	if len(res) == 0 {
		panic("empty response from rpc endpoint")
	}

	hexStr := res[0].Result

	n := new(big.Int)
	n.SetString(hexStr, 0)

	return int(n.Int64())
}

func (i *Indexer) Subscribe(address string) bool {
	table := (*i.AddressesDb).ReadOrCreate("eth_addresses")
	table.Set(address, address)

	return true
}

func (i *Indexer) GetTransactions(address string) []Transaction {
	table := (*i.TransactionsDb).ReadOrCreate("eth_transactions")
	txns := table.Get(address)

	if len(txns) == 0 {
		txns = []Transaction{}
	}

	return txns
}
