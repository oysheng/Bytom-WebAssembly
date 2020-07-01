package blockchain

import (
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
		in.AssetID = orig.AssetID()
		in.Amount = orig.Amount()
	} else {
		in.AssetID = *consensus.BTMAssetID
	}

	id := tx.Tx.InputIDs[i]
	in.InputID = id
	e := tx.Entries[id]
	switch e := e.(type) {
	case *bc.VetoInput:
		in.Type = "veto"
		in.ControlProgram = orig.ControlProgram()
		in.Address = getAddressFromControlProgram(in.ControlProgram, false)
		in.SpentOutputID = e.SpentOutputId
		arguments := orig.Arguments()
		for _, arg := range arguments {
			in.WitnessArguments = append(in.WitnessArguments, arg)
		}

	case *bc.CrossChainInput:
		in.Type = "cross_chain_in"
		in.ControlProgram = orig.ControlProgram()
		in.Address = getAddressFromControlProgram(in.ControlProgram, true)
		in.SpentOutputID = e.MainchainOutputId
		arguments := orig.Arguments()
		for _, arg := range arguments {
			in.WitnessArguments = append(in.WitnessArguments, arg)
		}

	case *bc.Spend:
		in.Type = "spend"
		in.ControlProgram = orig.ControlProgram()
		in.Address = getAddressFromControlProgram(in.ControlProgram, false)
		in.SpentOutputID = e.SpentOutputId
		arguments := orig.Arguments()
		for _, arg := range arguments {
			in.WitnessArguments = append(in.WitnessArguments, arg)
		}

	case *bc.Coinbase:
		in.Type = "coinbase"
		in.Arbitrary = e.Arbitrary
	}
	return in
}

// BuildAnnotatedOutput build the annotated output.
func BuildAnnotatedOutput(tx *types.Tx, idx int) *AnnotatedOutput {
	orig := tx.Outputs[idx]
	outid := tx.OutputID(idx)
	out := &AnnotatedOutput{
		OutputID:       *outid,
		Position:       idx,
		AssetID:        *orig.AssetAmount().AssetId,
		Amount:         orig.AssetAmount().Amount,
		ControlProgram: orig.ControlProgram(),
	}

	var isMainchainAddress bool
	switch e := tx.Entries[*outid].(type) {
	case *bc.IntraChainOutput:
		out.Type = "control"
		isMainchainAddress = false

	case *bc.CrossChainOutput:
		out.Type = "cross_chain_out"
		isMainchainAddress = true

	case *bc.VoteOutput:
		out.Type = "vote"
		out.Vote = e.Vote
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
