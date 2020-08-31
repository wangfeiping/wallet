package cosmos

import (
	"encoding/json"

	bip39 "github.com/cosmos/go-bip39"

	"github.com/wangfeiping/wallet/wallet/adapter/types"
)

const (
	// default number of words (12):
	// this generates a mnemonic directly from the number of words by reading system entropy.
	ENTROPY_SIZE = 128
)

var _ types.Adapter = (*AdapterCosmos)(nil)

type AdapterCosmos struct {
}

// CreateSeed returns mnemonics with bip39 to output 12-word list
func (a *AdapterCosmos) CreateSeed() string {
	entropy, err := bip39.NewEntropy(ENTROPY_SIZE)
	if err != nil {
		return err.Error()
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return err.Error()
	}
	//json output the seed
	var seed types.SeedOutput
	seed.Seed = mnemonic
	respbyte, _ := json.Marshal(seed)
	return string(respbyte)
}
