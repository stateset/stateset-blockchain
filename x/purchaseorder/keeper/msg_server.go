package keeper

import (
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
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
	k.zkpKeeper.VerifyProof(ctx, purchaseorder)

	// Add a DID to represent the Purchase Order in the Cosmosverse DID:STATESET:PO:123
	k.didKeeper.AddDID(ctx, purchaseorderhash)

	// Mint a NFT that represents the Purchase Order DID and Value of the PO
	k.bankKeeper.MintCoins(ctx, types.ModuleName, did)

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

func (server msgServer) CancelPurchaseOrder(goCtx context.Context, msg *types.MsgCancelPurchaseOrder) (*types.MsgCancelPurchaseOrderResponse, error) {
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
	purchaseorder.PurchaseOrderStatus = "cancelled"

	k.bankKeeper.BurnCoins(ctx, types.ModuleName, did)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtPurchaseOrderCancelled,
			sdk.NewAttribute(types.AttributeKeyPurchaseOrderId, strconv.FormatUint(msg.PurchaseOrderId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCancelPurchaseOrderResponse{}, nil
}


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

	// Draw Liquidity
	k.liquidityKeeper.DrawLiquidityFromCoins(coins, amount)

	// Create a Loan
	k.loanKeeper.CreateLoan(borrower, financer, collateral, amount)

	// Create a Stablecoin Covered Put Option
	k.optionKeeper.CreateStablecoinCoveredPutOption(asset, premium, expiration)

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