package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Markets []Market `json:"markets"`
	Params Params      `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Markets: nil,
		Params:  DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes market state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, market := range data.Markets {
		keeper.setMarket(ctx, market)
	}
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Markets: keeper.Markets(ctx),
		Params:      keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MinMarketLength < 1 {
		return fmt.Errorf("Param: MinMarketLength, must have a positive value")
	}

	if data.Params.MaxMarketLength < 1 {
		return fmt.Errorf("Param: MaxMarketLength, must have a positive value")
	}

	return nil
}