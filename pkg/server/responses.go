package server

import "github.com/PatricioNapoli/ethix/pkg/indexer"

type BlockResponse struct {
	Block int `json:"block"`
}

type SubscribeResponse struct {
	Subscribed bool `json:"subscribed"`
}

type TransactionsResponse struct {
	Transactions []indexer.Transaction `json:"transactions"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
