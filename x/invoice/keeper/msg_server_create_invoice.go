package keeper

import (
	"context"
	"fmt"

	"github.com/stateset/stateset-blockchain/x/invoice/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


// Create Invoice
func (server msgServer) CreateInvoice(goCtx context.Context, msg *types.MsgCreateInvoice) (*types.MsgCreateInvoiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	invoiceId, err := server.keeper.CreateInvoice(ctx, sender, msg.InvoiceParams, msg.InvoiceAssets)
	if err != nil {
		return nil, err
	}

	invoice, found := k.GetInvoice(ctx, msg.Id)
	invoice.InvoiceStatus = "created"

	// Verify the Value of the Invoice from existing system
	k.zkpKeeper.VerifyProof(ctx, invoice)
	
	// Add a DID to represent the Invoice in the Cosmosverse DID:STATESET:INV:123
	k.didKeeper.AddDID(ctx, invoicehash)
	
	// Mint a NFT that represents the Invoice DID and Value of the Invoice
	k.nftKeeper.MintCoins(ctx, did)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtInvoiceCreated,
			sdk.NewAttribute(types.AttributeKeyInvoiceId, strconv.FormatUint(purchaseOrderId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCreateInvoiceResponse{}, nil
}