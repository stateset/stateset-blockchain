package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stateset/stateset-blockchain/x/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AgreementAll(c context.Context, req *types.QueryAllAgreementRequest) (*types.QueryAllAgreementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var agreements []*types.Agreement
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	agreementStore := prefix.NewStore(store, types.KeyPrefix(types.AgreementKey))

	pageRes, err := query.Paginate(agreementStore, req.Pagination, func(key []byte, value []byte) error {
		var agreement types.Agreement
		if err := k.cdc.UnmarshalBinaryBare(value, &agreement); err != nil {
			return err
		}

		agreements = append(agreements, &agreement)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAgreementResponse{Agreement: agreements, Pagination: pageRes}, nil
}

func (k Keeper) Agreement(c context.Context, req *types.QueryGetAgreementRequest) (*types.QueryGetAgreementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var agreement types.Agreement
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AgreementKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.AgreementKey + req.Id)), &agreement)

	return &types.QueryGetAgreementResponse{Agreement: &agreement}, nil
}
