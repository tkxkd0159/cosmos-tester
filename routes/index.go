package routes

import sdktypes "github.com/line/lbm-sdk/types"

type MsgRouter map[string]func(m sdktypes.Msg) bool

func SetMsgRouter() MsgRouter {
	r := make(MsgRouter)
	r.registerBank()
	return r
}
