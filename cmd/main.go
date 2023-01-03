package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"

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

		txMsgData := new(sdk.TxMsgData)
		err := proto.Unmarshal(txres.GetData(), txMsgData)
		if err != nil {
			panic(err)
		}
		for j, m := range txMsgData.GetData() {
			fmt.Printf("  Msg[%d]: %s\n", j, m.GetMsgType())
		}

		for _, evt := range txres.Events {
			fmt.Println("\tEvent Type: ", evt.Type)
			for _, attr := range evt.GetAttributes() {
				fmt.Printf("\t\t%s\n", attr.String())
			}
		}
		fmt.Println()
	}
}
