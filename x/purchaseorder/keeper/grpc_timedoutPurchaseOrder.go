package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TimedoutPurchaseOrder(c context.Context, req *types.QueryGetTimedoutPurchaseOrderRequest) (*types.QueryGetTimedoutPurchaseOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var timedoutPurchaseOrder types.TimedoutPurchaseOrder
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasTimedoutPurchaseOrder(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetTimedoutPurchaseOrderIDBytes(req.Id)), &timedoutPurchaseOrder)

	return &types.QueryGetTimedoutPurchaseOrderResponse{TimedoutPurchaseOrder: &timedoutPurchaseOrder}, nil
}

func (k Keeper) TimedoutPurchaseOrders(c context.Context, req *types.QueryAllTimedoutPurchaseOrderRequest) (*types.QueryAllTimedoutPurchaseOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var timedoutPurchaseOrders []*types.TimedoutPurchaseOrder
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	timedoutPurchaseOrderStore := prefix.NewStore(store, types.KeyPrefix(types.TimedoutPurchaseOrderKey))

	pageRes, err := query.Paginate(timedoutPurchaseOrderStore, req.Pagination, func(key []byte, value []byte) error {
		var timedoutPurchaseOrder types.TimedoutPurchaseOrder
		if err := k.cdc.UnmarshalBinaryBare(value, &timedoutPurchaseOrder); err != nil {
			return err
		}

		timedoutPurchaseOrders = append(timedoutPurchaseOrders, &timedoutPurchaseOrder)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTimedoutPurchaseOrderResponse{TimedoutPurchaseOrder: timedoutPurchaseOrders, Pagination: pageRes}, nil
}
