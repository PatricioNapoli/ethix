package indexer

type Block struct {
	Hash         string        `json:"hash"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}
