package cosmo

import (
	"fmt"

	sdktypes "github.com/line/lbm-sdk/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
)

const BankMsgSend = "/cosmos.bank.v1beta1.MsgSend"

func FilterMsgs(m sdktypes.Msg) any {
	switch sdktypes.MsgTypeURL(m) {
	case BankMsgSend:
		res := m.(*banktypes.MsgSend)
		return res
	default:
		fmt.Printf("None\n")
	}
	return nil
}
