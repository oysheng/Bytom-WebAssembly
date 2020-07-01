package blockchain

import (
	"encoding/hex"

	"github.com/bytom-community/wasm/vapor/common"
	"github.com/bytom-community/wasm/vapor/consensus"
	"github.com/bytom-community/wasm/vapor/consensus/segwit"
	"github.com/bytom-community/wasm/vapor/protocol/bc"
	"github.com/bytom-community/wasm/vapor/protocol/bc/types"
)

// BuildAnnotatedInput build the annotated input.
func BuildAnnotatedInput(tx *types.Tx, i int) *AnnotatedInput {
	orig := tx.Inputs[i]
	in := &AnnotatedInput{}
	if orig.InputType() != types.CoinbaseInputType {
		assetID := orig.AssetID()
		in.AssetID = assetID.String()
		in.Amount = int64(orig.Amount())
		signData := tx.SigHash(uint32(i))
		in.SignData = signData.String()
		if vetoInput, ok := orig.TypedInput.(*types.VetoInput); ok {
			in.Vote = hex.EncodeToString(vetoInput.Vote)
		}
	} else {
		in.AssetID = consensus.BTMAssetID.String()
	}

	id := tx.Tx.InputIDs[i]
	in.InputID = id.String()
	e := tx.Entries[id]
	switch e := e.(type) {
	case *bc.VetoInput:
		in.Type = "veto"
		controlProgram := orig.ControlProgram()
		in.ControlProgram = hex.EncodeToString(controlProgram)
		in.Address = getAddressFromControlProgram(controlProgram, false)
		in.SpentOutputID = e.SpentOutputId.String()
		arguments := orig.Arguments()
		for _, arg := range arguments {
			in.WitnessArguments = append(in.WitnessArguments, hex.EncodeToString(arg))
		}

	case *bc.CrossChainInput:
		in.Type = "cross_chain_in"
		controlProgram := orig.ControlProgram()
		in.ControlProgram = hex.EncodeToString(controlProgram)
		in.Address = getAddressFromControlProgram(controlProgram, true)
		in.SpentOutputID = e.MainchainOutputId.String()
		arguments := orig.Arguments()
		for _, arg := range arguments {
			in.WitnessArguments = append(in.WitnessArguments, hex.EncodeToString(arg))
		}

	case *bc.Spend:
		in.Type = "spend"
		controlProgram := orig.ControlProgram()
		in.ControlProgram = hex.EncodeToString(controlProgram)
		in.Address = getAddressFromControlProgram(controlProgram, false)
		in.SpentOutputID = e.SpentOutputId.String()
		arguments := orig.Arguments()
		for _, arg := range arguments {
			in.WitnessArguments = append(in.WitnessArguments, hex.EncodeToString(arg))
		}

	case *bc.Coinbase:
		in.Type = "coinbase"
		in.Arbitrary = hex.EncodeToString(e.Arbitrary)
	}
	return in
}

// BuildAnnotatedOutput build the annotated output.
func BuildAnnotatedOutput(tx *types.Tx, idx int) *AnnotatedOutput {
	orig := tx.Outputs[idx]
	outputID := tx.OutputID(idx)
	out := &AnnotatedOutput{
		OutputID:       outputID.String(),
		Position:       idx,
		AssetID:        orig.AssetAmount().AssetId.String(),
		Amount:         int64(orig.AssetAmount().Amount),
		ControlProgram: hex.EncodeToString(orig.ControlProgram()),
	}

	var isMainchainAddress bool
	switch e := tx.Entries[*outputID].(type) {
	case *bc.IntraChainOutput:
		out.Type = "control"
		isMainchainAddress = false

	case *bc.CrossChainOutput:
		out.Type = "cross_chain_out"
		isMainchainAddress = true

	case *bc.VoteOutput:
		out.Type = "vote"
		out.Vote = hex.EncodeToString(e.Vote)
		isMainchainAddress = false
	}

	out.Address = getAddressFromControlProgram(orig.ControlProgram(), isMainchainAddress)
	return out
}

func getAddressFromControlProgram(prog []byte, isMainchain bool) string {
	netParams := &consensus.MainNetParams
	if isMainchain {
		netParams = consensus.BytomMainNetParams(&consensus.MainNetParams)
	}

	if segwit.IsP2WPKHScript(prog) {
		if pubHash, err := segwit.GetHashFromStandardProg(prog); err == nil {
			return buildP2PKHAddress(pubHash, netParams)
		}
	} else if segwit.IsP2WSHScript(prog) {
		if scriptHash, err := segwit.GetHashFromStandardProg(prog); err == nil {
			return buildP2SHAddress(scriptHash, netParams)
		}
	}
	return ""
}

func buildP2PKHAddress(pubHash []byte, netParams *consensus.Params) string {
	address, err := common.NewAddressWitnessPubKeyHash(pubHash, netParams)
	if err != nil {
		return ""
	}
	return address.EncodeAddress()
}

func buildP2SHAddress(scriptHash []byte, netParams *consensus.Params) string {
	address, err := common.NewAddressWitnessScriptHash(scriptHash, netParams)
	if err != nil {
		return ""
	}
	return address.EncodeAddress()
}
