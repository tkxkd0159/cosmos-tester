package cosmo

import (
	lbm "github.com/line/lbm-sdk/simapp"
	lbmparams "github.com/line/lbm-sdk/simapp/params"
)

func LbmEnc() lbmparams.EncodingConfig {
	return lbm.MakeTestEncodingConfig()
}
