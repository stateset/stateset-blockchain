package purchaseorder

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgNewPurchaseOrder:
			return handleMsgNewPurchaseOrder(ctx, keeper, msg)
		
		case types.MsgEditPurchaseOrder:
			return handleMsgEditPurchaseOrder(ctx, keeper, msg)
		
		case types.MsgCompletePurchaseOrder:
			return handleMsgCompletePurchaseOrder(ctx, keeper, msg)

		case types.MsgCancelPurchaseOrder:
			return handleMsgCancelPurchaseOrder(ctx, keeper, msg)

		case types.MsgFinancePurchaseOrder:
			return handleMsgFinancePurchaseOrder(ctx, handlerkeeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized purchase order message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgNewPurchaseOrder(ctx sdk.Context,  keeper Keeper, msg types.MsgNewPurchaseOrder) sdk.Result {
	order, err := types.NewPurchaseOrder(msg.Source, msg.Destination, msg.Owner, ctx.BlockTime(), msg.ClientOrderId)
	if err != nil {
		return err.Result()
	}


	return keeper.NewPurchaseOrderSingle(ctx, agreement)
}

func handleMsgEditPurchaseOrder(ctx sdk.Context,  keeper Keeper, msg types.MsgEditPurchaseOrder) sdk.Result {

	return keeper.EditPurchaseOrder(ctx, msg.Owner, msg.PurchaseOrderId)
}

func handleMsgCompletePurchaseOrder(ctx sdk.Context, keeper Keeper, msg types.MsgCompletePurchaseOrder) sdk.Result {

	return keeper.CompletePurchaseOrder(ctx, msg.Owner, msg.PurchaseOrderId)
}

func handleMsgCancelPurchaseOrder(ctx sdk.Context,  keeper Keeper, msg types.MsgCancelPurchaseOrder) sdk.Result {

	return keeper.MsgCancelPurchaseOrder(ctx, msg.Owner, msg.PurchaseOrderId)
}

func handleMsgFinancePurchaseOrder(ctx sdk.Context,  keeper Keeper, msg types.MsgFinancePurchaseOrder) sdk.Result {

	return keeper.FinancePurchaseOrder(ctx, msg.Owner, msg.PurchaseOrderId)(ctx, msg.Owner, msg.PurchaseOrderId)
}