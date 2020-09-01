package cosmos

import (
	// "github.com/cosmos/cosmos-sdk/client/keys"
	// crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/spf13/viper"

	"github.com/wangfeiping/wallet/wallet/adapter/types"
)

const (
	// default number of words (12):
	// this generates a mnemonic directly from the number of words by reading system entropy.
	defaultEntropySize = 128
	defaultBIP39pass   = ""
	defaultDenomName   = "uatom"
)

var _ types.Adapter = (*AdapterCosmos)(nil)

type AdapterCosmos struct {
}

// CreateSeed returns mnemonics with bip39 to output 12-word list
func (a *AdapterCosmos) CreateSeed() (string, error) {
	entropy, err := bip39.NewEntropy(defaultEntropySize)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// CreateAccount returns the json string for the created account.
func (a *AdapterCosmos) CreateAccount(rootDir, name,
	password, seed string) (ko types.KeyOutput, err error) {
	viper.Set("home", rootDir)
	kb := keyring.NewInMemory()

	//check out the input
	if name == "" {
		err = types.ErrMissingName
		return
	}
	if password == "" {
		err = types.ErrMissingPassword
		return
	}
	// check if already exists
	infos, err := kb.List()
	for _, info := range infos {
		if info.GetName() == name {
			err = types.ErrKeyNameConflict
			return
		}
	}

	//create account
	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString("secp256k1", keyringAlgos)
	if err != nil {
		return
	}
	if seed == "" {
		// algo := crkeys.SigningAlgo("secp256k1")
		name := "inmemorykey"
		var mnem string
		_, mnem, err = kb.NewMnemonic(name, keyring.English,
			defaultBIP39pass, algo)
		if err != nil {
			return
		}
		seed = mnem
	}

	// info, err1 := kb.NewAccount(name, seed, defaultBIP39pass, password, 0, 0)
	info, err := kb.NewAccount(name, seed, defaultBIP39pass, password, algo)
	if err != nil {
		return
	}

	// keyOutput, err2 := crkeys.Bech32KeyOutput(info)
	keyOutput, err := keyring.Bech32KeyOutput(info)
	if err != nil {
		return
	}

	//add new field denom for the coin name
	ko.Name = keyOutput.Name
	ko.Type = keyOutput.Type
	ko.Address = keyOutput.Address
	ko.PubKey = keyOutput.PubKey
	ko.Seed = seed
	ko.Denom = defaultDenomName
	return
}
