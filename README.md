
<p align="center">
    <img alt="Ethix" src="assets/img/ethix.jpg" width="400px"/>
</p>

<div align="center">

  <a style="margin-right:15px" href="#"><img src="https://forthebadge.com/images/badges/made-with-go.svg" alt="Made with Go"/></a>
  <a style="margin-right:15px" href="#"><img src="https://forthebadge.com/images/badges/powered-by-black-magic.svg" alt="Made with Go"/></a>
  <a href="https://www.paradigm.xyz/2020/08/ethereum-is-a-dark-forest"><img src="assets/img/dark-forest.svg" alt="Ethereum is a dark forest"/></a>


  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-brightgreen.svg" alt="License MIT"/></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/go-1.18-blue.svg" alt="Go 1.18"/></a>
  <a href="https://github.com/PatricioNapoli/ethix/actions/workflows/build.yml"><img src="https://github.com/PatricioNapoli/ethix/actions/workflows/build.yml/badge.svg" alt="build"/></a>

</div>


# Ethix

## Overview

ETHereum IndeXer, is a Go microservice that handles ethereum transaction indexing given an address. Only stores new transactions starting from the address subscription onwards, not historical. Meant to be used as a notifications service for new transactions, but still uses storage from the initial subscription starting point.

## Prerequisites

If you are not using docker, these are required:  

make  
go 1.18+

## Build & Run

### Make

#### Build & Run

`make go`  

Alternatively, you may run `make build` and `make run` separately.  
Or run the scripts in `scripts/`.  

#### Testing

`make test`  

### Docker

#### Building

`docker build -t patricionapoli/ethix .`  

#### Testing

`docker build -f t.Dockerfile -t patricionapoli/ethix .`

#### Running 

`docker run -p 8080:8080 patricionapoli/ethix`   

#### Running release

`docker run -p 8080:8080 ghcr.io/patricionapoli/ethix:master`

## Usage

Port used is 8080. You may use Curl for testing.

### POST /block

Gets current ETH mainnet block number.

```bash
$ curl "http://localhost:8080/block"

{"block":16243174}
```

### POST /subscribe?addr=address

Subscribes an ETH address from the query string, marking it for transaction search and parsing.

```bash
$ curl -v "http://localhost:8080/subscribe?addr=0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5"

{"subscribed":true}
```

### POST /transactions?addr=address

Gets all transactions from address from the subscription onwards.

```bash
$ curl "http://localhost:8080/transactions?addr=0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5"

{
	"transactions": [{
		"gas": "0x565f",
		"hash": "0x5d3968719230cf870132f438b8d6c74dceaebb28656ac7f33e82776d9938efbf",
		"input": "0x",
		"gasPrice": "0x2e69597f0",
		"maxFeePerGas": "0x2e69597f0",
		"maxPriorityFeePerGas": "0x0",
		"chainId": "0x1",
		"transactionIndex": "0x124",
		"value": "0x17e17b7d26e2684",
		"blockHash": "0xb0308003e06319cd18edf9c25b09904e6a5ee456f3375014b9cbcce884006702",
		"blockNumber": "0xf7d967",
		"accessList": [],
		"type": "0x2",
		"from": "0xdafea492d9c6733ae3d56b7ed1adb60692c98bc5",
		"to": "0x388c818ca8b9251b393131c08a736a67ccb19297",
		"nonce": "0x237f6",
		"v": "0x1",
		"r": "0xbeaf0ba413859895a93b0ab2f0a7dd2f8caea301ad4f718c3025414208dc5e52",
		"s": "0x4c6ef8e4deadd9d0c7689b93b005ce5e83573f749715c205a199fa216d3bfe2"
	}]
}
```
