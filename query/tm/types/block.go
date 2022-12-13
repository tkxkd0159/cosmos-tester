package types

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	types1 "github.com/tendermint/tendermint/proto/tendermint/types"
)

type ResultBlockResults struct {
	Height                int64                     `json:"height,string"`
	TxsResults            []*abci.ResponseDeliverTx `json:"txs_results"`
	BeginBlockEvents      []abci.Event              `json:"begin_block_events"`
	EndBlockEvents        []abci.Event              `json:"end_block_events"`
	ValidatorUpdates      []abci.ValidatorUpdate    `json:"validator_updates"`
	ConsensusParamUpdates *ConsensusParams          `json:"consensus_param_updates"`
}

type ConsensusParams struct {
	Block     *BlockParams            `json:"block,omitempty"`
	Evidence  *EvidenceParams         `json:"evidence,omitempty"`
	Validator *types1.ValidatorParams `json:"validator,omitempty"`
	Version   *VersionParams          `json:"version,omitempty"`
}

type BlockParams struct {
	// Note: must be greater than 0
	MaxBytes int64 `json:"max_bytes,omitempty,string"`
	// Note: must be greater or equal to -1
	MaxGas int64 `json:"max_gas,omitempty,string"`
}

type EvidenceParams struct {
	// Max age of evidence, in blocks.
	// The basic formula for calculating this is: MaxAgeDuration / {average block time}.
	MaxAgeNumBlocks int64 `json:"max_age_num_blocks,omitempty,string"`
	// Max age of evidence, in time.
	MaxAgeDuration time.Duration `json:"max_age_duration,string"`
	// This sets the maximum size of total evidence in bytes that can be committed in a single block.
	// and should fall comfortably under the max block bytes.
	// Default is 1048576 or 1MB
	MaxBytes int64 `json:"max_bytes,omitempty,string"`
}

// VersionParams contains the ABCI application version.
type VersionParams struct {
	AppVersion uint64 `json:"app_version,omitempty,string"`
}
