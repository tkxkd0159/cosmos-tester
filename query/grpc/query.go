package grpc

import (
	"context"
	"log"
	"strconv"

	tmtypes "github.com/line/lbm-sdk/client/grpc/tmservice"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
	"google.golang.org/grpc/metadata"
)

const HeaderHeight = "x-cosmos-block-height"

func (q *Querier) AllBalances(addr string, height ...int64) *banktypes.QueryAllBalancesResponse {
	bq := banktypes.NewQueryClient(q)
	var ctx context.Context
	if height != nil {
		header := metadata.New(map[string]string{HeaderHeight: strconv.FormatInt(height[0], 10)})
		ctx = metadata.NewOutgoingContext(context.Background(), header)
	} else {
		ctx = context.Background()
	}
	res, err := bq.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: addr})
	if err != nil {
		log.Println(err)
	}
	return res
}

func (q *Querier) BlockByHeight(height int64) *tmtypes.GetBlockByHeightResponse {
	tmq := tmtypes.NewServiceClient(q)
	res, err := tmq.GetBlockByHeight(context.Background(), &tmtypes.GetBlockByHeightRequest{Height: height})
	if err != nil {
		log.Println(err)
	}
	return res
}

func (q *Querier) LatestBlock() *tmtypes.GetLatestBlockResponse {
	tmq := tmtypes.NewServiceClient(q)
	res, err := tmq.GetLatestBlock(context.Background(), &tmtypes.GetLatestBlockRequest{})
	if err != nil {
		log.Println(err)
	}
	return res
}
