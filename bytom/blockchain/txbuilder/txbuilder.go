// Package txbuilder builds a Chain Protocol transaction from
// a list of actions.
package txbuilder

import (
	"context"

	"github.com/bytom-community/wasm/bytom/errors"
)

// errors
var (
	//ErrBadTxInputIdx means unsigned tx input
	ErrBadTxInputIdx = errors.New("unsigned tx missing input")
)

// Sign will try to sign all the witness
func Sign(ctx context.Context, tpl *Template, auth string, signFn SignFunc) error {
	for i, sigInst := range tpl.SigningInstructions {
		for j, wc := range sigInst.WitnessComponents {
			switch sw := wc.(type) {
			case *SignatureWitness:
				err := sw.sign(ctx, tpl, uint32(i), auth, signFn)
				if err != nil {
					return errors.WithDetailf(err, "adding signature(s) to signature witness component %d of input %d", j, i)
				}
			case *RawTxSigWitness:
				err := sw.sign(ctx, tpl, uint32(i), auth, signFn)
				if err != nil {
					return errors.WithDetailf(err, "adding signature(s) to raw-signature witness component %d of input %d", j, i)
				}
			}
		}
	}
	return materializeWitnesses(tpl)
}
