package keeper

import (
	"github.com/stateset/stateset-blockchain/x/invoice/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

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
	k.nftKeeper.MintNFT(ctx, did)

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

// Complete Invoice
func (server msgServer) CompleteInvoice(goCtx context.Context, msg *types.MsgCompleteInvoice) (*types.MsgCompleteInvoiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = server.keeper.CompleteInvoice(ctx, sender, msg.InvoiceId, msg.amount)
	if err != nil {
		return nil, err
	}

	invoice, found := k.GetInvoice(ctx, msg.Id)
	invoice.InvoiceStatus = "completed"

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtInvoiceCompleted,
			sdk.NewAttribute(types.AttributeKeyInvoiceId, strconv.FormatUint(msg.InvoiceId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCompleteInvoiceResponse{}, nil
}

// Factor Invoice
func (server msgServer) FactorInvoice(goCtx context.Context, msg *types.MsgFactorInvoice) (*types.MsgFactorInvoiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = server.keeper.FactorInvoice(ctx, sender, msg.InvoiceId, msg.Amount)
	if err != nil {
		return nil, err
	}

	invoice, found := k.GetInvoice(ctx, msg.Id)
	invoice.InvoiceStatus = "factored"

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtInvoiceFactored,
			sdk.NewAttribute(types.AttributeKeyInvoiceId, strconv.FormatUint(msg.InvoiceId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgFinanceInvoiceResponse{}, nil
}