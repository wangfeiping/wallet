package types

import (
	"fmt"
)

type SeedOutput struct {
	Seed  string `json:"seed"`
	Error string `json:"error"`
}

type KeyOutput struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Address      string `json:"address"`
	PubKey       string `json:"pub_key"`
	PrivKeyArmor string `json:"priv_key_armor,omitempty"`
	Seed         string `json:"seed,omitempty"`
	Denom        string `json:"denom"`
	Error        string `json:"error"`
}

type Adapter interface {
	// CreateSeed returns mnemonics with bip39 to output 12-word list
	// CreateSeed() (string, error)

	// CreateAccount returns the json string for the created account.
	CreateAccount(rootDir, name, password, seed string) (KeyOutput, error)
}

var (
	// flags
	FlagHome     = "home"
	FlagBitSize  = "bit.size"
	FlagName     = "name"
	FlagPasswd   = "passwd"
	FlagMnemonic = "mnemonic"

	// errors on account creation
	ErrKeyNameConflict = fmt.Errorf("acount with the name already exists")
	ErrMissingName     = fmt.Errorf("you have to specify a name for the locally stored account")
	ErrMissingPassword = fmt.Errorf("you have to specify a password for the locally stored account")
	ErrMissingSeed     = fmt.Errorf("you have to specify seed for key recover")
)
