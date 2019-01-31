package base

import (
	"encoding/json"
	"github.com/bytom-community/wasm/errors"
	"syscall/js"

	"github.com/bytom-community/wasm/common"
	"github.com/bytom-community/wasm/consensus"
	chainjson "github.com/bytom-community/wasm/encoding/json"
	"github.com/bytom-community/wasm/protocol/vm"
	"github.com/bytom-community/wasm/protocol/vm/vmutil"
	"github.com/bytom-community/wasm/sdk/lib"
)

// ContractArgument for smart contract
type ContractArgument struct {
	Type    string          `json:"type"`
	RawData json.RawMessage `json:"raw_data"`
}

// DataArgument is the other argument for run contract
type DataArgument struct {
	Value chainjson.HexBytes `json:"value"`
}

// StrArgument is the string argument for run contract
type StrArgument struct {
	Value string `json:"value"`
}

// IntegerArgument is the integer argument for run contract
type IntegerArgument struct {
	Value int64 `json:"value"`
}

// BoolArgument is the boolean argument for run contract
type BoolArgument struct {
	Value bool `json:"value"`
}

// AddressArgument is the address argument for run contract
type AddressArgument struct {
	Value string `json:"value"`
}

func ConvertContractArg(arg ContractArgument) (chainjson.HexBytes, error) {
	resultData := chainjson.HexBytes{}
	switch arg.Type {
	case "data":
		data := &DataArgument{}
		if err := json.Unmarshal(arg.RawData, data); err != nil {
			return nil, err
		}
		resultData = data.Value

	case "string":
		data := &StrArgument{}
		if err := json.Unmarshal(arg.RawData, data); err != nil {
			return nil, err
		}
		resultData = []byte(data.Value)

	case "integer":
		data := &IntegerArgument{}
		if err := json.Unmarshal(arg.RawData, data); err != nil {
			return nil, err
		}
		resultData = vm.Int64Bytes(data.Value)

	case "boolean":
		data := &BoolArgument{}
		if err := json.Unmarshal(arg.RawData, data); err != nil {
			return nil, err
		}
		resultData = vm.BoolBytes(data.Value)

	case "address":
		data := &AddressArgument{}
		if err := json.Unmarshal(arg.RawData, data); err != nil {
			return nil, err
		}

		address, err := common.DecodeAddress(data.Value, &consensus.ActiveNetParams)
		if err != nil {
			return nil, err
		}

		redeemContract := address.ScriptAddress()
		program := []byte{}
		switch address.(type) {
		case *common.AddressWitnessPubKeyHash:
			program, err = vmutil.P2WPKHProgram(redeemContract)
		case *common.AddressWitnessScriptHash:
			program, err = vmutil.P2WSHProgram(redeemContract)
		default:
			return nil, errors.New("bad address type")
		}
		resultData = program

	default:
		return nil, errors.New("bad argument type")
	}

	return resultData, nil
}

// ConvertArgument convert arguments
func ConvertArgument(args []js.Value) {
	defer lib.EndFunc(args[1]) //end func call
	typ := args[0].Get("type").String()
	if lib.IsEmpty(typ) {
		args[1].Set("error", "type empty")
		return
	}

	rawData := json.RawMessage{}
	err := json.Unmarshal([]byte(args[0].Get("raw_data").String()), &rawData)

	arg := ContractArgument{
		Type:    typ,
		RawData: rawData,
	}

	data, err := ConvertContractArg(arg)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	args[1].Set("data", string(data))
}
