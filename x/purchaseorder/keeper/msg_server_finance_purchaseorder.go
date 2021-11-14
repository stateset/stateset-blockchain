package keeper

import (
	"context"
	"fmt"

	"github.com/stateset/stateset-blockchain/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (server msgServer) FinancePurchaseOrder(goCtx context.Context, msg *types.MsgFinancePurchaseOrder) (*types.MsgFinancePurchaseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = server.keeper.FinancePurchaseOrder(ctx, sender, msg.PurchaseOrderId, msg.Amount)
	if err != nil {
		return nil, err
	}

	borrower, _ := sdk.AccAddressFromBech32(msg.Creator)

	financer, _ := sdk.AccAddressFromBech32(msg.Financer)

	collateral, _ := sdk.ParseCoins(msg.Collateral) // NFT of the PO

	amount, _ := sdk.ParseCoins(msg.Amount) // Fees

	// Send the NFT to the Loan Module to be Collateralized
	k.bankKeeper.SendCoinsFromModuleToModule(ctx, purchaseorder, loan)

	// Create a Loan
	k.loanKeeper.CreateLoanRequest(borrower, financer, collateral, amount)

	// Create a Stablecoin Covered Put Option
	//k.optionKeeper.CreateStablecoinCoveredPutOption(asset, premium, expiration)

	purchaseorder, found := k.GetPurchaseOrder(ctx, msg.Id)
	purchaseorder.PurchaseOrderStatus = "financed"

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtPurchaseOrderFinanced,
			sdk.NewAttribute(types.AttributeKeyPurchaseOrderId, strconv.FormatUint(msg.PurchaseOrderId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgFinancePurchaseOrderResponse{}, nil
}