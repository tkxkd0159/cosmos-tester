package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type User struct {
	Priv cryptotypes.PrivKey
	Pub  cryptotypes.PubKey
	Addr sdk.AccAddress
}
