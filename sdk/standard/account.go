package standard

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/bytom-community/wasm/bytom/account"
	"github.com/bytom-community/wasm/bytom/blockchain/signers"
	"github.com/bytom-community/wasm/bytom/blockchain/txbuilder"
	"github.com/bytom-community/wasm/bytom/common"
	"github.com/bytom-community/wasm/bytom/consensus"
	"github.com/bytom-community/wasm/bytom/crypto"
	"github.com/bytom-community/wasm/bytom/crypto/ed25519/chainkd"
	"github.com/bytom-community/wasm/bytom/crypto/sha3pool"
	"github.com/bytom-community/wasm/bytom/protocol/vm/vmutil"
	"github.com/bytom-community/wasm/sdk/lib"
)

//CreateAccount create account
func CreateAccount(this js.Value, args []js.Value) interface{} {
	defer lib.EndFunc(args[1])
	var (
		alias     string
		quorum    int
		rootXPub  string
		nextIndex uint64 // account next index like 1 2 3 ...
	)

	alias = args[0].Get("alias").String()
	quorum = args[0].Get("quorum").Int()
	rootXPub = args[0].Get("rootXPub").String()
	nextIndex = uint64(args[0].Get("nextIndex").Int())

	var XPubs []chainkd.XPub
	xpub := new(chainkd.XPub)
	xpub.UnmarshalText([]byte(rootXPub))
	XPubs = append(XPubs, *xpub)

	normalizedAlias := strings.ToLower(strings.TrimSpace(alias))

	signer, err := signers.Create("account", XPubs, quorum, nextIndex)
	id := signers.IDGenerate()
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}

	acc := &account.Account{Signer: signer, ID: id, Alias: normalizedAlias}
	rawAccount, err := json.Marshal(acc)
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}
	args[1].Set("data", string(rawAccount))
	return nil
}

//CreateAccountReceiver create address by account
func CreateAccountReceiver(this js.Value, args []js.Value) interface{} {
	defer lib.EndFunc(args[1])
	var (
		acc       account.Account
		err       error
		nextIndex uint64
		cp        *account.CtrlProgram
	)
	err = json.Unmarshal([]byte(args[0].Get("account").String()), &acc)
	nextIndex = uint64(args[0].Get("nextIndex").Int())

	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}
	if len(acc.XPubs) == 1 {
		cp, err = createP2PKH(&acc, false, nextIndex)
	} else {
		cp, err = createP2SH(&acc, false, nextIndex)
	}
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}

	res, err := controlPrograms(cp)
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}
	data, _ := json.Marshal(res)
	args[1].Set("db", string(data)) //insert web IndexedDB
	tx := txbuilder.Receiver{
		ControlProgram: cp.ControlProgram,
		Address:        cp.Address,
	}
	txj, _ := json.Marshal(tx)
	args[1].Set("data", string(txj))
	return nil
}

func createP2PKH(acc *account.Account, change bool, nextIndex uint64) (*account.CtrlProgram, error) {
	path := signers.Path(acc.Signer, signers.AccountKeySpace, nextIndex)
	derivedXPubs := chainkd.DeriveXPubs(acc.XPubs, path)
	derivedPK := derivedXPubs[0].PublicKey()
	pubHash := crypto.Ripemd160(derivedPK)

	address, err := common.NewAddressWitnessPubKeyHash(pubHash, &consensus.ActiveNetParams)
	if err != nil {
		return nil, err
	}

	control, err := vmutil.P2WPKHProgram([]byte(pubHash))
	if err != nil {
		return nil, err
	}

	return &account.CtrlProgram{
		AccountID:      acc.ID,
		Address:        address.EncodeAddress(),
		KeyIndex:       nextIndex,
		ControlProgram: control,
		Change:         change,
	}, nil
}

func createP2SH(acc *account.Account, change bool, nextIndex uint64) (*account.CtrlProgram, error) {
	path := signers.Path(acc.Signer, signers.AccountKeySpace, nextIndex)
	derivedXPubs := chainkd.DeriveXPubs(acc.XPubs, path)
	derivedPKs := chainkd.XPubKeys(derivedXPubs)
	signScript, err := vmutil.P2SPMultiSigProgram(derivedPKs, acc.Quorum)
	if err != nil {
		return nil, err
	}
	scriptHash := crypto.Sha256(signScript)

	address, err := common.NewAddressWitnessScriptHash(scriptHash, &consensus.ActiveNetParams)
	if err != nil {
		return nil, err
	}

	control, err := vmutil.P2WSHProgram(scriptHash)
	if err != nil {
		return nil, err
	}

	return &account.CtrlProgram{
		AccountID:      acc.ID,
		Address:        address.EncodeAddress(),
		KeyIndex:       nextIndex,
		ControlProgram: control,
		Change:         change,
	}, nil
}

func controlPrograms(progs ...*account.CtrlProgram) (map[string]string, error) {
	var hash common.Hash
	res := make(map[string]string)
	for _, prog := range progs {
		accountCP, err := json.Marshal(prog)
		if err != nil {
			return nil, err
		}

		sha3pool.Sum256(hash[:], prog.ControlProgram)
		res[account.ContractKeyHexString(hash)] = string(accountCP)
	}
	return res, nil
}

type PubKeyResp struct {
	XPub        string   `json:"xpub"`
	Pubkey      string   `json:"pubkey"`
	DerivedPath []string `json:"derived_path"`
}

// CreatePubkey create pubkey
func CreatePubkey(this js.Value, args []js.Value) interface{} {
	defer lib.EndFunc(args[1]) //end func call
	xpubStr := args[0].Get("xpub").String()
	if lib.IsEmpty(xpubStr) || len(xpubStr) != 128 {
		args[1].Set("error", fmt.Sprintf("invalid xpub:", xpubStr))
		return nil
	}

	xpubByte, err := hex.DecodeString(xpubStr)
	if err != nil {
		args[1].Set("error", "decode xpub")
		return nil
	}
	var xpub chainkd.XPub
	copy(xpub[:], xpubByte)
	pubkey := xpub.PublicKey()

	seed := args[0].Get("seed").Int()
	if seed <= 0 {
		args[1].Set("error", "invalid seed with not positive integer")
		return nil
	}

	derivedPath := []string{}
	path := signers.Path(&signers.Signer{KeyIndex: uint64(1)}, signers.AccountKeySpace, uint64(seed))
	for _, p := range path {
		derivedPath = append(derivedPath, hex.EncodeToString(p))
	}

	res := PubKeyResp{
		XPub:        xpubStr,
		Pubkey:      hex.EncodeToString(pubkey),
		DerivedPath: derivedPath,
	}
	rawPubkeyResp, err := json.Marshal(res)
	if err != nil {
		args[1].Set("error", err.Error())
		return nil
	}

	args[1].Set("data", string(rawPubkeyResp))
	return nil
}
