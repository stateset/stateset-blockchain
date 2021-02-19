package handler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	mkeeper "github.com/stateset/stateset-blockchain/x/agreement/keeper"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/keeper"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)


func NewHandler(keeper keeper.IKeeper, mkeeper mkeeper.IKeeper) sdk.Handler {
	ms := NewMsgServerImpl(keeper, mkeeper)

	func NewHandler(keeper keeper.IKeeper, mkeeper mkeeper.IKeeper) sdk.Handler {
		ms := NewMsgServerImpl(keeper, mkeeper)
		case *types.MsgCreatePurchaseOrder:
			res, err := ms.CreatePurchaseOrder(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		
		case *types.MsgEditPurchaseOrder:
			res, err := ms.EditPurchaseOrder(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCompletePurchaseOrder:
			res, err := ms.CancelPurchaseOrder(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		
		case *types.MsgCancelPurchaseOrder:
			res, err := ms.CancelPurchaseOrder(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgFinancePurchaseOrder:
			res, err := ms.FinancePurchaseOrder(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized purchaseorder message type: %T", msg)
		}
	}
}