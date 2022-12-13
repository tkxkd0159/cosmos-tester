package main

import (
	"fmt"

	"cosmo/query/tm"
)

const (
	RESTADDR = "https://cosmos-mainnet-rpc.allthatnode.com:26657"
)

func main() {
	client := tm.NewQuerier(RESTADDR)
	targetHeight := 13229976
	res, err := client.BlockResults(targetHeight)
	if err != nil {
		panic(err)
	}
	for _, txres := range res.TxsResults {
		fmt.Println(txres.Code)
		for _, evt := range txres.Events {
			fmt.Println(evt.String())
		}
	}
}
