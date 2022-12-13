package routes

import (
	"fmt"

	sdktypes "github.com/line/lbm-sdk/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
)

func (mr MsgRouter) registerBank() {
	mr[sdktypes.MsgTypeURL(&banktypes.MsgSend{})] = handleMsgSend
}

func handleMsgSend(sm sdktypes.Msg) bool {
	m, _ := sm.(*banktypes.MsgSend)
	fmt.Println("Handle: ", m.FromAddress, m.ToAddress, m.Amount)
	return true
}
