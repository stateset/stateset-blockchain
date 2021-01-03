package app

import (
	"encoding/json"
)

// GenesisState is the state of the blockchain.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for stateset.
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}