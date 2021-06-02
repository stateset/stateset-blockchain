package invoice

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stateset/stateset-blockchain/x/invoice/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InvoiceAll(c context.Context, req *types.QueryAllInvoiceRequest) (*types.QueryAllInvoiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var invoices []*types.Invoice
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	invoiceStore := prefix.NewStore(store, types.KeyPrefix(types.InvoiceKey))

	pageRes, err := query.Paginate(invoiceStore, req.Pagination, func(key []byte, value []byte) error {
		var invoice types.Invoice
		if err := k.cdc.UnmarshalBinaryBare(value, &invoice); err != nil {
			return err
		}

		invoices = append(invoices, &invoice)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInvoiceResponse{Invoice: invoices, Pagination: pageRes}, nil
}

func (k Keeper) Invoice(c context.Context, req *types.QueryGetInvoiceRequest) (*types.QueryGetInvoiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var invoice types.Invoice
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InvoiceKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.InvoiceKey + req.Id)), &invoice)

	return &types.QueryGetInvoiceResponse{Invoice: &invoice}, nil
}
