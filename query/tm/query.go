package tm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	jsonrpc "github.com/tendermint/tendermint/rpc/jsonrpc/types"

	ctypes "cosmo/query/tm/types"
)

func (q *Querier) ParseJSONRPC(r *http.Response) (*jsonrpc.RPCResponse, error) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	res := new(jsonrpc.RPCResponse)
	err = res.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (q *Querier) BlockResults(height int) (*ctypes.ResultBlockResults, error) {
	query := fmt.Sprintf("/block_results?height=%d", height)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, q.URL+query, nil)
	resp, err := q.Do(req)
	if err != nil {
		return nil, err
	}
	res, err := q.ParseJSONRPC(resp)
	if err != nil {
		return nil, err
	}
	blockRes := new(ctypes.ResultBlockResults)
	err = json.Unmarshal(res.Result, blockRes)
	if err != nil {
		return nil, err
	}
	return blockRes, nil
}
