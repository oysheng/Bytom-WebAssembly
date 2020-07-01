package side

import (
	"encoding/json"
	"syscall/js"

	"github.com/bytom-community/wasm/sdk/lib"
	"github.com/bytom-community/wasm/vapor/blockchain"
	"github.com/bytom-community/wasm/vapor/common/arithmetic"
	"github.com/bytom-community/wasm/vapor/protocol/bc/types"
)

// decodeVaporRawTx decode vapor raw transaction
func decodeVaporRawTx(rawVaporTx string) (*blockchain.AnnotatedRawTx, error) {
	var tx *types.Tx
	if err := tx.UnmarshalText([]byte(rawVaporTx)); err != nil {
		return nil, err
	}

	annotatedTx := &blockchain.AnnotatedRawTx{
		ID:        tx.ID,
		Version:   tx.Version,
		Size:      tx.SerializedSize,
		TimeRange: tx.TimeRange,
		Inputs:    []*blockchain.AnnotatedInput{},
		Outputs:   []*blockchain.AnnotatedOutput{},
	}

	for i := range tx.Inputs {
		annotatedTx.Inputs = append(annotatedTx.Inputs, blockchain.BuildAnnotatedInput(tx, i))
	}
	for i := range tx.Outputs {
		annotatedTx.Outputs = append(annotatedTx.Outputs, blockchain.BuildAnnotatedOutput(tx, i))
	}

	annotatedTx.Fee, _ = arithmetic.CalculateTxFee(tx)
	return annotatedTx, nil
}

// DecodeVaporRawTx decode vapor raw transaction
func DecodeVaporRawTx(this js.Value, args []js.Value) interface{} {
	defer lib.EndFunc(args[1]) //end func call
	rawTx := args[0].Get("raw_transaction").String()
	if lib.IsEmpty(rawTx) {
		args[1].Set("error", "raw_transaction empty")
		return nil
	}

	vaporTx, err := decodeVaporRawTx(rawTx)
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}

	rawVaporTx, err := json.Marshal(vaporTx)
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}

	args[1].Set("data", string(rawVaporTx))
	return nil
}
