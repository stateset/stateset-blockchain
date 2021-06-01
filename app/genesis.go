package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GenesisState is the state of the stateset blockchain.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the blockchain.
func NewDefaultGenesisState(cdc codec.JSONMarshaler) GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}