package cosmo

import (
	"fmt"

	sdktypes "github.com/line/lbm-sdk/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
)

func FilterMsgs(m sdktypes.Msg) any {
	anym := m.(any)
	switch anym.(type) {
	case *banktypes.MsgSend:
		fmt.Println("Route here", anym)
	default:
		fmt.Println("Not Machted: ", anym)
	}
	return nil
}

func DecodeTx(tx sdktypes.Tx) {
	for _, m := range tx.GetMsgs() {
		FilterMsgs(m)
		fmt.Println(m.String(), sdktypes.MsgTypeURL(m))
		for i := 0; i < len(m.GetSigners()); i++ {
			fmt.Printf("Signer [%d]: %s\n", i, m.GetSigners()[i])
		}
	}
}
