package keeper

import (
	"context"
	"fmt"

	"github.com/stateset/stateset-blockchain/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (server msgServer) CreatePurchaseOrder(goCtx context.Context, msg *types.MsgCreatePurchaseOrder) (*types.MsgCreatePurchaseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	poolId, err := server.keeper.CreatePurchaseOrder(ctx, sender, msg.PurchaseOrderParams, msg.PurchaseOrderAssets)
	if err != nil {
		return nil, err
	}

	purchaseorder, found := k.GetPurchaseOrder(ctx, msg.Id)
	purchaseorder.PurchaseOrderStatus = "created"

	// Verify the Value of the Purchase Order from existing system
	// k.zkpKeeper.VerifyProof(ctx, purchaseorder)

	// Add a DID to represent the Purchase Order in the Cosmosverse DID:STATESET:PO:123
	k.didKeeper.AddDID(ctx, purchaseorderhash)

	// Mint a NFT that represents the Purchase Order DID and Value of the PO
	k.bankKeeper.MintCoins(ctx, purchaseorder, 1)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtPurchaseOrderCreated,
			sdk.NewAttribute(types.AttributeKeyPurchaseOrderId, strconv.FormatUint(purchaseOrderId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCreatePurchaseOrderResponse{}, nil
}