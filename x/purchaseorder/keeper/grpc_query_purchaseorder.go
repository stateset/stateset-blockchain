package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stateset/stateset-blockchain/x/stateset-blockchain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PurchaseorderAll(c context.Context, req *types.QueryAllPurchaseorderRequest) (*types.QueryAllPurchaseorderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var purchaseorders []*types.Purchaseorder
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	purchaseorderStore := prefix.NewStore(store, types.KeyPrefix(types.PurchaseorderKey))

	pageRes, err := query.Paginate(purchaseorderStore, req.Pagination, func(key []byte, value []byte) error {
		var purchaseorder types.Purchaseorder
		if err := k.cdc.UnmarshalBinaryBare(value, &purchaseorder); err != nil {
			return err
		}

		purchaseorders = append(purchaseorders, &purchaseorder)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPurchaseorderResponse{Purchaseorder: purchaseorders, Pagination: pageRes}, nil
}

func (k Keeper) Purchaseorder(c context.Context, req *types.QueryGetPurchaseorderRequest) (*types.QueryGetPurchaseorderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var purchaseorder types.Purchaseorder
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PurchaseorderKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.PurchaseorderKey + req.Id)), &purchaseorder)

	return &types.QueryGetPurchaseorderResponse{Purchaseorder: &purchaseorder}, nil
}
