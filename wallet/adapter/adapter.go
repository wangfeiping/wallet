package adapter

import (
	"encoding/json"

	"github.com/wangfeiping/wallet/wallet/adapter/cosmos"
	"github.com/wangfeiping/wallet/wallet/adapter/types"
)

var cosmosAdapter types.Adapter

func init() {
	cosmosAdapter = &cosmos.AdapterCosmos{}
}

func CreateSeed() string {
	mnem, err := cosmos.CreateSeed()
	// json output the seed
	var seed types.SeedOutput
	if err != nil {
		seed.Error = err.Error()
	} else {
		seed.Seed = mnem
	}
	bytes, _ := json.Marshal(seed)
	return string(bytes)
}

func CreateAccount(rootDir, name, passwd, seed string) string {
	acc, err := cosmosAdapter.CreateAccount(rootDir, name, passwd, seed)
	if err != nil {
		acc.Error = err.Error()
	}
	bytes, _ := json.Marshal(acc)
	return string(bytes)
}

// RecoverKey returns the account info that recovered with name, password and mnemonic input
func RecoverKey(rootDir, name, passwd, seed string) string {
	acc, err := cosmosAdapter.CreateAccount(rootDir, name, passwd, seed)
	if err != nil {
		acc.Error = err.Error()
	}
	acc.Seed = ""
	bytes, _ := json.Marshal(acc)
	return string(bytes)
}
