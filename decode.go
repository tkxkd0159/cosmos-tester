package cosmo

import (
	"fmt"

	lbm "github.com/line/lbm-sdk/simapp"
	lbmparams "github.com/line/lbm-sdk/simapp/params"
	sdktypes "github.com/line/lbm-sdk/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"

	"cosmo/routes"
)

func LbmEnc() lbmparams.EncodingConfig {
	return lbm.MakeTestEncodingConfig()
}

func DecodeTx(tx sdktypes.Tx) {
	r := routes.SetMsgRouter()
	for _, m := range tx.GetMsgs() {
		RouteMsgsByType(m)

		if handle := r[sdktypes.MsgTypeURL(m)]; handle != nil {
			handle(m)
		}

		fmt.Println(m.String(), sdktypes.MsgTypeURL(m))
		for i := 0; i < len(m.GetSigners()); i++ {
			fmt.Printf("Signer [%d]: %s\n", i, m.GetSigners()[i])
		}
	}
}

func RouteMsgsByType(m sdktypes.Msg) any {
	anym := m.(any)
	switch anym.(type) {
	case *banktypes.MsgSend:
		fmt.Println("Route here", anym)
	default:
		fmt.Println("Not Machted: ", anym)
	}
	return nil
}
