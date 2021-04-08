package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stateset/stateset-blockchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k Keeper) Agreements(c context.Context, req *types.QueryAgreementsRequest) (*types.QueryAgreementsResponse, error) {
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

func (k Keeper) Agreement(ctx context.Context, req *types.QueryAgreementRequest) (*types.QueryAgreementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var agreement types.Agreement
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AgreementKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.AgreementKey + req.Id)), &agreement)

	return &types.QueryAgreementResponse{Agreement: &agreement}, nil
}


func (k Keeper) AgreementParams(ctx context.Context, req *types.QueryAgreementParamsRequest) (*types.QueryAgreementParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	agreement, err := k.GetAgreement(sdkCtx, req.AgreementId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAgreementParamsResponse{
		Params: agreement.GetAgreementParams(),
	}, nil
}

