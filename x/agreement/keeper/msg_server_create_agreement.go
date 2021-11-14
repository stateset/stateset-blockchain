package keeper

import (
	"context"
	"fmt"

	"github.com/stateset/stateset-blockchain/x/agreement/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (server msgServer) CreateAgreement(goCtx context.Context, msg *types.MsgCreateAgreement) (*types.MsgCreateAgreementResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	poolId, err := server.keeper.CreateAgreement(ctx, sender, msg.AgreementParams, msg.AgremenetAssets)
	if err != nil {
		return nil, err
	}

	agreement, found := k.GetAgreement(ctx, msg.Id)
	agreement.AgreementStatus = "created"

	// Verify the Value of the Agreement from existing system
	k.zkpKeeper.VerifyProof(ctx, agreement)
	
	// Add a DID to represent the Agreement in the Cosmosverse DID:STATESET:AGREEMENT:123
	k.didKeeper.AddDID(ctx, agreementhash)

	// Mint a NFT that represents the Agreement DID and Value of the Agreement
	k.nftKeeper.MintNFT(ctx, did)

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