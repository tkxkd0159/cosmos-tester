# cosmos-tester

## Data Aggregation Strategy
아래의 옵션들 중 선택
1. TMRPC의 `block_results`, gRPC의 `BlockByHeight`의 결과를 동시에 받아온 후 `block_results`에서 성공 tx에 대한 인덱스 거르고 실제 tx 내용은 `BlockByHeight`에서 받아서 DB에 저장
2. Osmosis에서 pool 관련 event 리스트업 하고 `block_results`에서  `code: 0`인 success tx에 대해 내가 원하는 타겟 이벤트들 걸러서 DB에 관련 내용 저장

* `EmitTypedEvent`로 만들어진 이벤트의 경우 `type`이 <proto_package_name>.<message_name>.
* `message.action`, `message.sender`의 경우 `baseapp.runMsgs`에서 자동 생성. `message.module`의 경우 msgHandler에서 주입하지 않은 경우 여기서 주입.
```shell
# query="<type>.<attributes_key>='<attribute_value>'"
GET https://rpc-cosmoshub.blockapsis.com/tx_search?query="coin_spent.spender='cosmos1vqem22tk25epwzu0zcjdmd8famwezy3nl3namq'"
> https://ebony-rpc.finschia.io/tx_search?query=%22coin_spent.spender=%27tlink1ucxztrucmnxy8hhgmxte2knsz3gh7y5yzmj83s%27%22
> https://ebony-rpc.finschia.io/tx_search?query=%22timeout_packet.packet_src_port=%27transfer%27%22
> https://ebony-rpc.finschia.io/tx_search?query=%22message.module=%27bank%27%22
> https://ebony-rpc.finschia.io/tx_search?query=%22message.action=%27/cosmos.bank.v1beta1.MsgSend%27%22

total_count보고 pagination 가능.
> https://rpc.testnet.osmosis.zone/tx_search?query=%22message.action=%27/ibc.applications.transfer.v1.MsgTransfer%27%22&order_by=%22desc%22&per_page=100&page=1
아니면 아래처럼 특정 높이(5744465) 이후로만 이벤트 매칭되는 것 찾을 수도 있음.
> https://rpc.testnet.osmosis.zone/tx_search?query=%22tx.height%20%3E%205744465%20AND%20message.action=%27/ibc.applications.transfer.v1.MsgTransfer%27%22&order_by=%22asc%22&per_page=100

# Another query example
"tm.Event='Tx' AND message.action=/cosmos.bank.v1.MsgSend"
```

## Index TX
Tendermint에서 인덱싱을 진행함. `kv`와 `psql` 두가지 방법 선택 가능. `kv`의 경우 쿼리가 제한적이고 추후 없어질 예정. `psql`의 경우 `state/indexer/sink/psql/schema.sql`을 사전에 미리 테이블 세팅해야 함.  
`tx.height`, `tx.hash`, `block.height`의 경우 abci.Event 생성 시 Index를 true로 설정하지 않아도 자동으로 인덱싱됨. 나머지는 `Index: true` 설정해야 인덱싱 됨.

TODO: 과거 데이터를 어떻게 깔끔하게 검색할 수 있을까? 외부에서 인덱싱하고 이 내용을 체인에서 검증하는 방식으로 고민 필요

## Query raw store(state) (ABCI Query)
### 1. gRPC (/{ServiceName}/{MethodName})
[Invoke](https://github.com/cosmos/cosmos-sdk/blob/cdc329189b0fa11c13b36faf27d7f5bf96a19ff4/client/grpc_query.go#L32)를 통해 조회. `Invoke`에서 `BroadcastTxRequest`가 아닌 경우 query로 핸들링한다. Handler의 경우 `RegisterService(sd *grpc.ServiceDesc, handler interface{})`를 통해 app 시작 시 등록됨.
```go
	opts := rpcclient.ABCIQueryOptions{
		Height: queryHeight,
		Prove:  true|false,
	}

  result, err := node.ABCIQueryWithOptions(context.Background(), req.Path, req.Data, opts)
```
아래 Balance 쿼리를 예시로 들면 Invoke에서 Path는 `/cosmos.bank.v1beta1.Query/Balance`, Data는 req인 `QueryBalanceRequest`를 `[]byte`로 마샬링한 값이다. 
```go
func (c *queryClient) Balance(ctx context.Context, in *QueryBalanceRequest, opts ...grpc.CallOption) (*QueryBalanceResponse, error) {
	out := new(QueryBalanceResponse)
	err := c.cc.Invoke(ctx, "/cosmos.bank.v1beta1.Query/Balance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
```
해당값으로 node.ABCIQueryWithOptions를 날리면 응답으로 아래 데이터가 온다.
```go
type ResponseQuery struct {
	Code uint32
	Log       string
	Info      string
	Index     int64
	Key       []byte
	Value     []byte
	ProofOps  *crypto.ProofOps
	Height    int64
	Codespace string
}
```
여기서 Value 값을 reply인 `QueryBalanceResponse`로 언마샬링하면 원하는 쿼리 데이터를 얻음.
Consensus client로 바로 조회할 때는 `<rpcaddr>:26657/abci_query?path=_&data=_&height=_&prove=true`

### 2. /app Query
* `/app/simulate`를 Path로 txBytes를 Data로 넣어서 해당 tx의 실횅 결과를 `SimulationResponse`로 받을 수 있음.
* `/app/version`으로 Path 설정해서 baseapp.version 값을 받을 수 있음.

### 3. /store Query
기본적으로 latest height -1를 조회함. 즉, header height = data height + 1이고 latest height - 1이여야 proof를 가져올 수 있음. latest height 기준 값이 궁금하면 직접 그 높이를 설정해서 쿼리해야 함. Path는 `/store/<substore>/<path>` 형태로 들어감. 이 때 <substore>는 storeName이고 rootMultiStore에서 매칭되는 모듈의 store를 불러올 때 사용.
* `/store/<storeName>/key`로 쿼리할 경우 `req.Data`의 경우 이 데이터를 들고 있는 실제 key byte가 들어감. 이것은 추후 ABCI Query의 응답인 `res.Key` 값에도 들어감. `res.Value`에는 해당 store의 key(`req.Data`)에 저장된 값을 설정함. `res.ProofOps`는 해당 데이터에 대한 merkle proof를 불러와서 저장함. Prefix store의 경우에도 `<prefix_byte><key_byte>` 같은 방식으로 `Data`를 설정하면 됨.
* `/store/<storeName>/subspace`로 쿼리할 경우 해당 store에 저장된 모든 key, value 쌍을 iterator를 통해 전부 불러와서 반환함.

```go
func UpgradedClientKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s/%d/%s", KeyUpgradedIBCState, height, KeyUpgradedClient))
}

RequestQuery{
  Path: "/store/upgrade/key",
  Height: upgradedHeight,
  Data: "upgradedIBCState/<upgradedHeight>/upgradedClient",
  Prove: true
}

store.Set(key, value)
```

### 4. /p2p Query
`/p2p/<cmd>/<type>/<arg>` 형태로 구성됨. App단에 설정된 PeerFilter를 통해 peer를 <ip:port>나 node ID로 필터링할 정보를 얻기 위해 사용.
* `/p2p/filter/addr/<arg>`, `/p2p/filter/id/<arg>` :
```go
	addrQuery := abci.RequestQuery{
		Path: "/p2p/filter/addr/1.1.1.1:8000",
	}

	idQuery := abci.RequestQuery{
		Path: "/p2p/filter/id/testid",
	}
```
