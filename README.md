# cosmos-tester

## Data Aggregation Strategy
아래의 옵션들 중 선택
1. TMRPC의 `block_results`, gRPC의 `BlockByHeight`의 결과를 동시에 받아온 후 `block_results`에서 성공 tx에 대한 인덱스 거르고 실제 tx 내용은 `BlockByHeight`에서 받아서 DB에 저장
2. Osmosis에서 pool 관련 event 리스트업 하고 `block_results`에서  `code: 0`인 success tx에 대해 내가 원하는 타겟 이벤트들 걸러서 DB에 관련 내용 저장
```shell
# query="<type>.<attributes_key>='<attribute_value>'"
GET https://rpc-cosmoshub.blockapsis.com/tx_search?query="coin_spent.spender='cosmos1vqem22tk25epwzu0zcjdmd8famwezy3nl3namq'"
```

## Index TX
Tendermint에서 인덱싱을 진행함. `kv`와 `psql` 두가지 방법 선택 가능. `kv`의 경우 쿼리가 제한적이고 추후 없어질 예정. `psql`의 경우 `state/indexer/sink/psql/schema.sql`을 사전에 미리 테이블 세팅해야 함.  
`tx.height`, `tx.hash`, `block.height`의 경우 abci.Event 생성 시 Index를 true로 설정하지 않아도 자동으로 인덱싱됨. 나머지는 `Index: true` 설정해야 인덱싱 됨.
