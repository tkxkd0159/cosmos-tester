package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrant "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"google.golang.org/protobuf/proto"

	"cosmo/compress_test/lib"
	t "cosmo/compress_test/types"
)

const (
	LOAD         = 1000
	ChainID      = "test"
	CurrentBatch = 0
)

func main() {
	start := time.Now()
	batchByte := setupBatch(LOAD)
	fmt.Println("Serialized batch size:", len(batchByte), "bytes")
	batch := new(t.Batch)
	err := proto.Unmarshal(batchByte, batch)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Println("Setup took", elapsed.String())
}

func setupBatch(n int) []byte {
	users := make([]t.User, n)
	for i := 0; i < n; i++ {
		priv, pub, addr := testdata.KeyTestPubAddr()
		users[i] = t.User{Priv: priv, Pub: pub, Addr: addr}
	}

	encCfg := testutil.MakeTestEncodingConfig(
		auth.AppModuleBasic{}, authz.AppModuleBasic{}, bank.AppModuleBasic{}, distribution.AppModuleBasic{},
		evidence.AppModuleBasic{}, feegrant.AppModuleBasic{}, genutil.AppModuleBasic{}, gov.AppModuleBasic{},
		mint.AppModuleBasic{}, slashing.AppModuleBasic{}, staking.AppModuleBasic{}, upgrade.AppModuleBasic{}, vesting.AppModuleBasic{})
	txBuilder := encCfg.TxConfig.NewTxBuilder()

	txs := make([][]byte, n)
	blocks := make([]*t.MockBlock, LOAD)
	batches := make([]*t.Batch, 0)
	var err error
	for i, u := range users {
		txs[i], err = genTx(txBuilder, encCfg.TxConfig, u, users[(i+1)%len(users)]) // 1 ~ 5ms
		if err != nil {
			panic(err)
		}
		blocks[i] = &t.MockBlock{Txs: [][]byte{txs[i]}}
	}
	batches = append(batches, &t.Batch{Elements: blocks})

	batchByte, err := proto.Marshal(batches[CurrentBatch])
	if err != nil {
		panic(err)
	}
	return batchByte
}

func genTx(b client.TxBuilder, txCfg client.TxConfig, from, to t.User) ([]byte, error) {
	err := b.SetMsgs(banktypes.NewMsgSend(from.Addr, to.Addr, sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(rand.Int63n(100000))))))
	if err != nil {
		return nil, err
	}

	// Body
	b.SetMemo(lib.RandomWord())
	b.SetTimeoutHeight(rand.Uint64())

	// AuthInfo
	b.SetGasLimit(rand.Uint64())
	b.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(rand.Int63n(100000)))))

	// Signatures
	accNum := rand.Uint64()
	accSeq := rand.Uint64()
	signerData := authsigning.SignerData{
		ChainID:       ChainID, // need for sign
		AccountNumber: accNum,  // need for sign
		Sequence:      accSeq,
	}
	bytesToSign, err := txCfg.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signerData, b.GetTx())
	if err != nil {
		return nil, err
	}
	sigBytes, err := from.Priv.Sign(bytesToSign)
	if err != nil {
		return nil, err
	}
	sig := signing.SignatureV2{
		PubKey: from.Pub,
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: sigBytes,
		},
		Sequence: accSeq,
	}
	err = b.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	return txCfg.TxEncoder()(b.GetTx())
}
