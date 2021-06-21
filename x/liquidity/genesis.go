package liquidity

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stateset/stateset-blockchain/x/liquidity/keeper"
	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Pools []Pool `json:"pools"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Pools:  nil,
		Params: DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(cdc codec.JSONMarshaler) GenesisState { return NewGenesisState() }

// InitGenesis new liquidity genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	// validate logic on module.go/InitGenesis
	for _, record := range data.Pools {
		keeper.SetPool(ctx, &record)
	}
	keeper.InitGenesis(ctx, data)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return keeper.ExportGenesis(ctx)
}
