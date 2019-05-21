// +build !mini

package js

import (
	"syscall/js"

	"github.com/bytom-community/wasm/sdk/standard"

	"github.com/bytom-community/wasm/sdk/base"
)

//RegisterFunc Register js func
type RegisterFunc func(this js.Value, args []js.Value) interface{}

var funcs map[string]RegisterFunc

func init() {
	funcs = make(map[string]RegisterFunc)

	funcs["createKey"] = base.CreateKey
	funcs["resetKeyPassword"] = base.ResetKeyPassword
	funcs["createAccount"] = standard.CreateAccount
	funcs["createAccountReceiver"] = standard.CreateAccountReceiver
	funcs["signTransaction"] = standard.SignTransaction
	funcs["signTransaction1"] = base.SignTransaction1
	funcs["convertArgument"] = base.ConvertArgument
	funcs["createPubkey"] = standard.CreatePubkey
}

//Register Register func
func Register() {
	jsFuncVal := js.Global().Get("AllFunc")
	for k, v := range funcs {
		call := js.FuncOf(v)
		jsFuncVal.Set(k, call)
	}
	setPrintMessage := js.Global().Get("setFuncOver")
	setPrintMessage.Invoke()
}
