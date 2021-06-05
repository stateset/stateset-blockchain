package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/ibc/applications/ibcdex/types"
)

func (k msgServer) CancelSellOrder(goCtx context.Context, msg *types.MsgCancelSellOrder) (*types.MsgCancelSellOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCancelSellOrderResponse{}, nil
}
