package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/stateset/stateset-blockchain/types"
	"github.com/stateset/stateset-blockchain/x/wasm/internal/keeper"
)

// BeginBlocker handles softfork over param changes
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	if core.IsSoftforkHeight(ctx, 1) {
		params := k.GetParams(ctx)
		params.MaxContractMsgSize = 4096
		k.SetParams(ctx, params)
	}
}