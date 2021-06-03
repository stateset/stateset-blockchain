package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stateset/stateset-blockchain/x/agreement/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TimedoutAgreementAll(c context.Context, req *types.QueryAllTimedoutAgreementRequest) (*types.QueryAllTimedoutAgreementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var timedoutAgreements []*types.TimedoutAgreement
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	timedoutAgreementStore := prefix.NewStore(store, types.KeyPrefix(types.TimedoutAgreementKey))

	pageRes, err := query.Paginate(timedoutAgreementStore, req.Pagination, func(key []byte, value []byte) error {
		var timedoutAgreement types.TimedoutAgreement
		if err := k.cdc.UnmarshalBinaryBare(value, &timedoutAgreement); err != nil {
			return err
		}

		timedoutAgreements = append(timedoutAgreements, &timedoutAgreement)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTimedoutAgreementResponse{TimedoutAgreement: timedoutAgreements, Pagination: pageRes}, nil
}

func (k Keeper) TimedoutAgreement(c context.Context, req *types.QueryGetTimedoutAgreementRequest) (*types.QueryGetTimedoutAgreementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var timedoutAgreement types.TimedoutAgreement
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasTimedoutAgreement(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutAgreementKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetTimedoutAgreementIDBytes(req.Id)), &timedoutAgreement)

	return &types.QueryGetTimedoutAgreementResponse{TimedoutAgreement: &timedoutAgreement}, nil
}
