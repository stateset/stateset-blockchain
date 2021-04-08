package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GenesisState is the state of the blockchain.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONMarshaler) GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}