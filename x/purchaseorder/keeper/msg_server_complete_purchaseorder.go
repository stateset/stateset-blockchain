package keeper

import (
	"context"
	"fmt"

	"github.com/stateset/stateset-blockchain/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (server msgServer) CompletePurchaseOrder(goCtx context.Context, msg *types.MsgCompletePurchaseOrder) (*types.MsgCompletePurchaseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = server.keeper.CompletePurchaseOrder(ctx, sender, msg.PurchaseOrderId, msg.amount)
	if err != nil {
		return nil, err
	}

	purchaseorder, found := k.GetPurchaseOrder(ctx, msg.Id)
	purchaseorder.PurchaseOrderStatus = "completed"

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtPurchaseOrderCompleted,
			sdk.NewAttribute(types.AttributeKeyPurchaseOrderId, strconv.FormatUint(msg.PurchaseOrderId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCompletePurchaseOrderResponse{}, nil
}