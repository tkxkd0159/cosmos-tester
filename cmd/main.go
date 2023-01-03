package main

import (
	"fmt"

	"cosmo/query/tm"
)

const (
	RESTADDR = "https://rpc-cosmoshub.blockapsis.com"
)

func main() {
	client := tm.NewQuerier(RESTADDR)
	targetHeight := 13514226
	res, err := client.BlockResults(targetHeight)
	if err != nil {
		panic(err)
	}
	for i, txres := range res.TxsResults {
		// 0 is success
		fmt.Printf("TX[%d] Code: %d\n", i, txres.Code)
		for _, evt := range txres.Events {
			fmt.Println("Event Type: ", evt.Type)
			for _, attr := range evt.GetAttributes() {
				fmt.Printf("   %s\n", attr.String())
			}
		}
		fmt.Println()
	}
}
