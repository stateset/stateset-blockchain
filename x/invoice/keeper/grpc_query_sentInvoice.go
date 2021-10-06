package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stateset/stateset-blockchain/x/invoices/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SentInvoice(c context.Context, req *types.QueryGetSentInvoiceRequest) (*types.QueryGetSentInvoiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sentInvoice types.SentInvoice
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasSentInvoice(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentInvoiceKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetSentInvoiceIDBytes(req.Id)), &sentInvoice)

	return &types.QueryGetSentInvoiceResponse{SentInvoice: &sentInvoice}, nil
}

func (k Keeper) SentInvoices(c context.Context, req *types.QueryAllSentInvoiceRequest) (*types.QueryAllSentInvoiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sentInvoices []*types.SentInvoice
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sentInvoiceStore := prefix.NewStore(store, types.KeyPrefix(types.SentInvoiceKey))

	pageRes, err := query.Paginate(sentInvoiceStore, req.Pagination, func(key []byte, value []byte) error {
		var sentInvoice types.SentInvoice
		if err := k.cdc.UnmarshalBinaryBare(value, &sentInvoice); err != nil {
			return err
		}

		sentInvoices = append(sentInvoices, &sentInvoice)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSentInvoiceResponse{SentInvoice: sentInvoices, Pagination: pageRes}, nil
}
