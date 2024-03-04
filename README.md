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
