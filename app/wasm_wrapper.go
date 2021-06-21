package app

import (
	"github.com/cosmos/cosmos-sdk/server"
)

type WasmWrapper struct {
	Wasm wasm.Config `mapstructure:"wasm"`
}

func getWasmConfig() wasm.Config {
	wasmWrap := WasmWrapper{Wasm: wasm.DefaultWasmConfig()}
	ctx := server.NewDefaultContext()
	err := ctx.Viper.Unmarshal(&wasmWrap)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	return wasmWrap.Wasm
}
