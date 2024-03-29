package indexer

import "encoding/json"

type Transaction struct {
	Gas                  string        `json:"gas"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	GasPrice             string        `json:"gasPrice"`
	MaxFeePerGas         string        `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	ChainId              string        `json:"chainId"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	AccessList           []interface{} `json:"accessList"`
	Type                 string        `json:"type"`
	From                 string        `json:"from"`
	To                   string        `json:"to"`
	Nonce                string        `json:"nonce"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	type _Transaction Transaction
	test := &_Transaction{
		AccessList: []interface{}{},
	}

	err := json.Unmarshal(data, test)
	if err != nil {
		return err
	}

	*t = Transaction(*test)
	return nil
}
