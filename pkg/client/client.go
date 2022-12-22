package client

import (
	"bytes"
	"github.com/PatricioNapoli/ethix/pkg/utils"
	"io/ioutil"
	"net/http"
)

type RPC[T any] struct {
	Endpoint string
	Version  string
}

type RPCRequest struct {
	RPCMethod string        `json:"method"`
	Params    []interface{} `json:"params"`
	Id        int           `json:"id"`
	Version   string        `json:"jsonrpc"`
}

type RPCResponse[T any] struct {
	Result T   `json:"result"`
	Id     int `json:"id"`
}

func NewRpc[T any](endpoint string, version string) *RPC[T] {
	return &RPC[T]{
		Endpoint: endpoint,
		Version:  version,
	}
}

func (rpc *RPC[T]) PostRequests(requests []*RPCRequest) (error, []RPCResponse[T]) {
	for i, r := range requests {
		r.Id = i
		r.Version = rpc.Version
	}

	json, err := utils.ToJSON(requests)
	if err != nil {
		return err, []RPCResponse[T]{}
	}

	responseBody := bytes.NewBuffer(json)

	resp, err := http.Post(rpc.Endpoint, "application/json", responseBody)
	if err != nil {
		return err, []RPCResponse[T]{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, []RPCResponse[T]{}
	}

	rpcRes := []RPCResponse[T]{}
	err = utils.FromJSON(body, &rpcRes)
	if err != nil {
		return err, []RPCResponse[T]{}
	}

	return nil, rpcRes
}
