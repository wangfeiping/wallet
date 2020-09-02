package adapter

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/wangfeiping/wallet/wallet/adapter/cosmos"
	"github.com/wangfeiping/wallet/wallet/adapter/ethereum"
	"github.com/wangfeiping/wallet/wallet/adapter/types"
)

var cosmosAdapter types.Adapter
var ethereumAdapter types.Adapter

var KeysCdc *codec.LegacyAmino

func init() {
	KeysCdc = codec.New()
	cosmosAdapter = &cosmos.AdapterCosmos{}
	ethereumAdapter = &ethereum.AdapterEthereum{}
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

// CosmosCreateAccount returns the account info that created with name, password and mnemonic input
func CosmosCreateAccount(rootDir, name, passwd, seed string) string {
	acc, err := cosmosAdapter.CreateAccount(rootDir, name, passwd, seed)
	if err != nil {
		acc.Error = err.Error()
	}
	bytes, _ := json.Marshal(acc)
	return string(bytes)
}

// CosmosRecoverKey returns the account info that recovered with name, password and mnemonic input
func CosmosRecoverKey(rootDir, name, passwd, seed string) string {
	acc, err := cosmosAdapter.CreateAccount(rootDir, name, passwd, seed)
	if err != nil {
		acc.Error = err.Error()
	}
	acc.Seed = ""
	// bytes, _ := json.Marshal(acc)
	bytes, _ := KeysCdc.MarshalJSON(acc)
	return string(bytes)
}

// Ethereum part

// EthCreateAccount returns the account info that created with name, password and mnemonic input.
func EthCreateAccount(rootDir, name, passwd, seed string) string {
	acc, err := ethereumAdapter.CreateAccount(rootDir, name, passwd, seed)
	if err != nil {
		acc.Error = err.Error()
	}
	bytes, _ := json.Marshal(acc)
	return string(bytes)
}

func EthRecoverAccount(rootDir, name, passwd, seed string) string {
	acc, err := ethereumAdapter.CreateAccount(rootDir, name, passwd, seed)
	if err != nil {
		acc.Error = err.Error()
	}
	acc.Seed = ""
	bytes, _ := json.Marshal(acc)
	return string(bytes)
}
