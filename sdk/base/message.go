package base

import (
	"encoding/hex"
	"encoding/json"
	"syscall/js"

	"github.com/bytom-community/wasm/sdk/lib"
)

// RespSignMessage is the response of SignMessage
type RespSignMessage struct {
	Signature string `json:"signature"`
}

// SignMessage sign message
func SignMessage(this js.Value, args []js.Value) interface{} {
	defer lib.EndFunc(args[1])
	message := args[0].Get("message").String()
	password := args[0].Get("password").String()
	keyJSON := args[0].Get("key").String()
	if lib.IsEmpty(message) || lib.IsEmpty(password) || lib.IsEmpty(keyJSON) {
		args[1].Set("error", "args empty")
		return nil
	}

	signData, err := SignData(keyJSON, nil, []byte(message), password)
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}

	var ret RespSignMessage
	ret.Signature = hex.EncodeToString(signData)
	j, err := json.Marshal(ret)
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}
	args[1].Set("data", string(j))
	return nil
}
