package keeper

import (
	"context"
	"fmt"

	"github.com/stateset/stateset-blockchain/x/invoice/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)



// Cancel Invoice
func (server msgServer) CancelInvoice(goCtx context.Context, msg *types.MsgCancelInvoice) (*types.MsgCancelInvoiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = server.keeper.CancelInvoice(ctx, sender, msg.InvoiceId, msg.amount)
	if err != nil {
		return nil, err
	}

	// Burn a NFT that represents the Invoice DID and Value of the Invoice
	k.bankKeeper.BurnCoins(ctx, did)

	invoice, found := k.GetInvoice(ctx, msg.Id)
	invoice.InvoiceStatus = "cancelled"

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtInvoiceCanceld,
			sdk.NewAttribute(types.AttributeKeyInvoiceId, strconv.FormatUint(msg.InvoiceId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCancelInvoiceResponse{}, nil
}