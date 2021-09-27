package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PurchaseOrder(c context.Context, req *types.QueryGetPurchaseOrderRequest) (*types.QueryGetPurchaseOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var purchaseorder types.PurchaseOrder
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PurchaseOrderKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.PurchaseOrderKey + req.Id)), &purchaseorder)

	return &types.QueryGetPurchaseOrderResponse{PurchaseOrder: &purchaseorder}, nil
}

func (k Keeper) PurchaseOrders(c context.Context, req *types.QueryAllPurchaseOrderRequest) (*types.QueryAllPurchaseOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var purchaseorders []*types.PurchaseOrder
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	purchaseOrderStore := prefix.NewStore(store, types.KeyPrefix(types.PurchaseOrderKey))

	pageRes, err := query.Paginate(purchaseOrderStore, req.Pagination, func(key []byte, value []byte) error {
		var purchaseorder types.PurchaseOrder
		if err := k.cdc.UnmarshalBinaryBare(value, &purchaseOrder); err != nil {
			return err
		}

		purchaseorders = append(purchaseorders, &purchaseorder)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPurchaseOrderResponse{PurchaseOrder: purchaseorders, Pagination: pageRes}, nil
}
