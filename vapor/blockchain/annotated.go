package blockchain

import (
	"encoding/json"

	"github.com/bytom-community/wasm/vapor/crypto/ed25519/chainkd"
	chainjson "github.com/bytom-community/wasm/vapor/encoding/json"
	"github.com/bytom-community/wasm/vapor/protocol/bc"
)

//AnnotatedTx means an annotated transaction.
type AnnotatedTx struct {
	ID                     bc.Hash            `json:"tx_id"`
	Timestamp              uint64             `json:"block_time"`
	BlockID                bc.Hash            `json:"block_hash"`
	BlockHeight            uint64             `json:"block_height"`
	Position               uint32             `json:"block_index"`
	BlockTransactionsCount uint32             `json:"block_transactions_count,omitempty"`
	Inputs                 []*AnnotatedInput  `json:"inputs"`
	Outputs                []*AnnotatedOutput `json:"outputs"`
	StatusFail             bool               `json:"status_fail"`
	Size                   uint64             `json:"size"`
}

// AnnotatedRawTx means an annotated raw transaction.
type AnnotatedRawTx struct {
	ID        bc.Hash            `json:"tx_id"`
	Version   uint64             `json:"version"`
	Size      uint64             `json:"size"`
	TimeRange uint64             `json:"time_range"`
	Inputs    []*AnnotatedInput  `json:"inputs"`
	Outputs   []*AnnotatedOutput `json:"outputs"`
	Fee       uint64             `json:"fee"`
}

//AnnotatedInput means an annotated transaction input.
type AnnotatedInput struct {
	Type             string   `json:"type"`
	InputID          string   `json:"input_id"`
	AssetID          string   `json:"asset"`
	Amount           int64    `json:"amount"`
	ControlProgram   string   `json:"script,omitempty"`
	Address          string   `json:"address,omitempty"`
	IssuanceProgram  string   `json:"issuance_program,omitempty"`
	AssetDefinition  string   `json:"asset_definition,omitempty"`
	SpentOutputID    string   `json:"spent_output_id,omitempty"`
	Arbitrary        string   `json:"arbitrary,omitempty"`
	WitnessArguments []string `json:"arguments,omitempty"`
	Vote             string   `json:"vote,omitempty"`
	SignData         string   `json:"sign_data,omitempty"`
}

//AnnotatedOutput means an annotated transaction output.
type AnnotatedOutput struct {
	Type           string `json:"type"`
	OutputID       string `json:"utxo_id"`
	Position       int    `json:"position"`
	AssetID        string `json:"asset"`
	Amount         int64  `json:"amount"`
	ControlProgram string `json:"script"`
	Address        string `json:"address,omitempty"`
	Vote           string `json:"vote,omitempty"`
}

//AnnotatedAccount means an annotated account.
type AnnotatedAccount struct {
	ID         string         `json:"id"`
	Alias      string         `json:"alias,omitempty"`
	XPubs      []chainkd.XPub `json:"xpubs"`
	Quorum     int            `json:"quorum"`
	KeyIndex   uint64         `json:"key_index"`
	DeriveRule uint8          `json:"derive_rule"`
}

//AnnotatedAsset means an annotated asset.
type AnnotatedAsset struct {
	ID                bc.AssetID         `json:"id"`
	Alias             string             `json:"alias"`
	VMVersion         uint64             `json:"vm_version"`
	RawDefinitionByte chainjson.HexBytes `json:"raw_definition_byte"`
	Definition        *json.RawMessage   `json:"definition"`
}

//AnnotatedUTXO means an annotated utxo.
type AnnotatedUTXO struct {
	Alias               string `json:"account_alias"`
	OutputID            string `json:"id"`
	AssetID             string `json:"asset_id"`
	AssetAlias          string `json:"asset_alias"`
	Amount              uint64 `json:"amount"`
	AccountID           string `json:"account_id"`
	Address             string `json:"address"`
	ControlProgramIndex uint64 `json:"control_program_index"`
	Program             string `json:"program"`
	SourceID            string `json:"source_id"`
	SourcePos           uint64 `json:"source_pos"`
	ValidHeight         uint64 `json:"valid_height"`
	Change              bool   `json:"change"`
	DeriveRule          uint8  `json:"derive_rule"`
}
