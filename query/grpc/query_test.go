package grpc_test

import (
	"os"
	"testing"

	"cosmo"
)

func TestDecodeTx(t *testing.T) {
	// client := grpc.NewQuerier("127.0.0.1:9090")
	// res := client.AllBalances("link146asaycmtydq45kxc8evntqfgepagygelel00h", 10)
	// fmt.Println(res)
	//
	// res2 := client.BlockByHeight(1974)
	// d := res2.GetBlock().GetData()
	// err := os.WriteFile("./test/tx", d.GetTxs()[0], 0o700)
	// if err != nil {
	// 	panic(err)
	// }
	// dat0 := d.GetTxs()[0]
	// fmt.Println(string(dat0))

	dat, err := os.ReadFile("../../test/tx")
	if err != nil {
		panic(err)
	}

	Enc := cosmo.LbmEnc()
	tx, err := Enc.TxConfig.TxDecoder()(dat)
	if err != nil {
		panic(err)
	}
	cosmo.DecodeTx(tx)
}
