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

func (k Keeper) SentPurchaseOrder(c context.Context, req *types.QueryGetSentPurchaseOrderRequest) (*types.QueryGetSentPurchaseOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sentPurchaseOrder types.SentPurchaseOrder
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasSentPurchaseOrder(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetSentPurchaseOrderIDBytes(req.Id)), &sentPurchaseOrder)

	return &types.QueryGetSentPurchaseOrderResponse{SentPurchaseOrder: &sentPurchaseOrder}, nil
}

func (k Keeper) SentPurchaseOrders(c context.Context, req *types.QueryAllSentPurchaseOrderRequest) (*types.QueryAllSentPurchaseOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sentPurchaseOrders []*types.SentPurchaseOrder
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sentPurchaseOrderStore := prefix.NewStore(store, types.KeyPrefix(types.SentPurchaseOrderKey))

	pageRes, err := query.Paginate(sentPurchaseOrderStore, req.Pagination, func(key []byte, value []byte) error {
		var sentPurchaseOrder types.SentPurchaseOrder
		if err := k.cdc.UnmarshalBinaryBare(value, &sentPurchaseOrder); err != nil {
			return err
		}

		sentPurchaseOrders = append(sentPurchaseOrders, &sentPurchaseOrder)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSentPurchaseOrderResponse{SentPurchaseOrder: sentPurchaseOrders, Pagination: pageRes}, nil
}
