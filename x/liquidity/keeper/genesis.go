package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

// new liquidity genesis
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := k.ValidateGenesis(ctx, genState); err != nil {
		panic(err)
	}
	k.SetParams(ctx, genState.Params)
	for _, record := range genState.PoolRecords {
		if err := k.ValidatePoolRecord(ctx, &record); err != nil {
			panic(err)
		}
		k.SetPoolRecord(ctx, &record)
	}
	// TODO: reset heights variables when init or export
}

// ValidateGenesis performs genesis state validation for the liquidity module.
func (k Keeper) ValidateGenesis(ctx sdk.Context, genState types.GenesisState) error {
	if err := genState.Params.Validate(); err != nil {
		return err
	}
	for _, record := range genState.PoolRecords {
		k.SetPoolRecord(ctx, &record)
		if err := k.ValidatePoolRecord(ctx, &record); err != nil {
			return err
		}
	}
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	var poolRecords []types.PoolRecord

	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		record, found := k.GetPoolRecord(ctx, pool)
		if found {
			poolRecords = append(poolRecords, *record)
		}
	}
	if len(poolRecords) == 0 {
		poolRecords = []types.PoolRecord{}
	}
	return types.NewGenesisState(params, poolRecords)
}