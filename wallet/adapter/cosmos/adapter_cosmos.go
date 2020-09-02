package cosmos

import (
	// "github.com/cosmos/cosmos-sdk/client/keys"
	// crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
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

// CreateSeed returns mnemonics with bip39 to output 12-word list
func CreateSeed() (string, error) {
	bitSize := viper.GetInt(types.FlagBitSize)
	if bitSize <= 0 {
		bitSize = defaultEntropySize
	}
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

var _ types.Adapter = (*AdapterCosmos)(nil)

type AdapterCosmos struct {
}

// CreateAccount returns the account info that created with name, password and mnemonic input.
func (a *AdapterCosmos) CreateAccount(rootDir, name,
	password, seed string) (ko types.KeyOutput, err error) {
	viper.Set(types.FlagHome, rootDir)
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
	// algo, err := keyring.NewSigningAlgoFromString("ed25519", keyringAlgos)
	if err != nil {
		return
	}
	// if seed == "" {
	// 	// algo := crkeys.SigningAlgo("secp256k1")
	// 	name := "inmemorykey"
	// 	var mnem string
	// 	_, mnem, err = kb.NewMnemonic(name, keyring.English,
	// 		defaultBIP39pass, algo)
	// 	if err != nil {
	// 		return
	// 	}
	// 	seed = mnem
	// }

	// info, err1 := kb.NewAccount(name, seed, defaultBIP39pass, password, 0, 0)

	// /github.com/cosmos/cosmos-sdk@v0.34.4-0.20200829041113-200e88ba075b/types/address.go
	// Atom in https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	// CoinType = 118
	hdPath := hd.CreateHDPath(118, 0, 0).String()
	var info keyring.Info
	info, err = kb.NewAccount(name, seed, defaultBIP39pass, hdPath, algo)
	if err != nil {
		return
	}

	// keyOutput, err2 := crkeys.Bech32KeyOutput(info)
	var keyOutput keyring.KeyOutput
	keyOutput, err = keyring.Bech32KeyOutput(info)
	if err != nil {
		return
	}

	//add new field denom for the coin name
	ko.Name = keyOutput.Name
	ko.Type = keyOutput.Type
	ko.Address = keyOutput.Address
	ko.PubKey = keyOutput.PubKey
	ko.Denom = defaultDenomName
	ko.Seed = seed

	// var armor string
	// armor, err = kb.ExportPrivKeyArmorByAddress(info.GetAddress(), password)
	// if err != nil {
	// 	return
	// }
	// ko.PrivKeyArmor = armor
	return
}
