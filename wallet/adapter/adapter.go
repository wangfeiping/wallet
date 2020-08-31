package adapter

import (
	"github.com/wangfeiping/wallet/wallet/adapter/cosmos"
	"github.com/wangfeiping/wallet/wallet/adapter/types"
)

var cosmosAdapter types.Adapter

func init() {
	cosmosAdapter = &cosmos.AdapterCosmos{}
}

func CreateSeed() string {
	return cosmosAdapter.CreateSeed()
}
