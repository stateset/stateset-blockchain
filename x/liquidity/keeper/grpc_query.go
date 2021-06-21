package keeper

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Make response of query liquidity pool
func (k Keeper) MakeQueryPoolResponse(ctx sdk.Context, pool types.Pool) (*types.QueryPoolResponse, error) {
	batch, found := k.GetPoolBatch(ctx, pool.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}

	return &types.QueryPoolResponse{Pool: pool,
		PoolMetadata: k.GetPoolMetaData(ctx, pool),
		PoolBatch:    batch}, nil
}

// Make response of query liquidity pools
func (k Keeper) MakeQueryPoolsResponse(ctx sdk.Context, pools types.Pools) (*[]types.QueryPoolResponse, error) {
	resp := make([]types.QueryPoolResponse, len(pools))
	for i, pool := range pools {
		batch, found := k.GetPoolBatch(ctx, pool.PoolId)
		if !found {
			return nil, types.ErrPoolBatchNotExists
		}
		meta := k.GetPoolMetaData(ctx, pool)
		res := types.QueryPoolResponse{
			Pool:         pool,
			PoolMetadata: meta,
			PoolBatch:    batch,
		}
		resp[i] = res
	}

	return &resp, nil
}

// read data from kvstore for response of query liquidity pool
func (k Keeper) Pool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	empty := &types.QueryPoolRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "liquidity pool %d doesn't exist", req.PoolId)
	}
	return k.MakeQueryPoolResponse(ctx, pool)
}

// read data from kvstore for response of query liquidity pools
func (k Keeper) Pools(c context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	empty := &types.QueryPoolsRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	poolStore := prefix.NewStore(store, types.PoolKeyPrefix)
	var pools types.Pools

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		pool, err := types.UnmarshalPool(k.cdc, value)
		if err != nil {
			return err
		}
		pools = append(pools, pool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err := k.MakeQueryPoolsResponse(ctx, pools)

	return &types.QueryPoolsResponse{*response, pageRes}, nil
}

// read data from kvstore for response of query liquidity pools batch
func (k Keeper) PoolsBatch(c context.Context, req *types.QueryPoolsBatchRequest) (*types.QueryPoolsBatchResponse, error) {
	empty := &types.QueryPoolsBatchRequest{}
	if req == nil || req == empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	batchStore := prefix.NewStore(store, types.PoolBatchKeyPrefix)
	var response []types.QueryPoolBatchResponse

	pageRes, err := query.Paginate(batchStore, req.Pagination, func(key []byte, value []byte) error {
		batch, err := types.UnmarshalPoolBatch(k.cdc, value)
		if err != nil {
			return err
		}
		res := &types.QueryPoolBatchResponse{
			Batch: batch,
		}
		response = append(response, *res)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsBatchResponse{response, pageRes}, nil
}

// read data from kvstore for response of query liquidity pools batc
func (k Keeper) PoolBatch(c context.Context, req *types.QueryPoolBatchRequest) (*types.QueryPoolBatchResponse, error) {
	empty := &types.QueryPoolBatchRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	batch, found := k.GetPoolBatch(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "liquidity pool batch %d doesn't exist", req.PoolId)
	}
	return &types.QueryPoolBatchResponse{Batch: batch}, nil
}

// read data from kvstore for response of query batch swap messages of the liquidity pool batch
func (k Keeper) PoolBatchSwapMsgs(c context.Context, req *types.QueryPoolBatchSwapMsgsRequest) (*types.QueryPoolBatchSwapMsgsResponse, error) {
	empty := &types.QueryPoolBatchSwapMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.PoolBatchSwapMsgIndexKeyPrefix)
	var msgs []types.BatchPoolSwapMsg

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalBatchPoolSwapMsg(k.cdc, value)
		if err != nil {
			return err
		}

		msgs = append(msgs, msg)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolBatchSwapMsgsResponse{
		Swaps:      msgs,
		Pagination: pageRes,
	}, nil
	return nil, nil
}

// read data from kvstore for response of query batch deposit messages of the liquidity pool batch
func (k Keeper) PoolBatchDepositMsgs(c context.Context, req *types.QueryPoolBatchDepositMsgsRequest) (*types.QueryPoolBatchDepositMsgsResponse, error) {
	empty := &types.QueryPoolBatchDepositMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.PoolBatchDepositMsgIndexKeyPrefix)
	var msgs []types.BatchPoolDepositMsg

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalBatchPoolDepositMsg(k.cdc, value)
		if err != nil {
			return err
		}

		msgs = append(msgs, msg)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolBatchDepositMsgsResponse{
		Deposits:   msgs,
		Pagination: pageRes,
	}, nil
}

// read data from kvstore for response of query batch withdraw messages of the liquidity pool batch
func (k Keeper) PoolBatchWithdrawMsgs(c context.Context, req *types.QueryPoolBatchWithdrawMsgsRequest) (*types.QueryPoolBatchWithdrawMsgsResponse, error) {
	empty := &types.QueryPoolBatchWithdrawMsgsRequest{}
	if req == nil || *req == *empty {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	msgStore := prefix.NewStore(store, types.PoolBatchWithdrawMsgIndexKeyPrefix)
	var msgs []types.BatchPoolWithdrawMsg

	pageRes, err := query.Paginate(msgStore, req.Pagination, func(key []byte, value []byte) error {
		msg, err := types.UnmarshalBatchPoolWithdrawMsg(k.cdc, value)
		if err != nil {
			return err
		}

		msgs = append(msgs, msg)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolBatchWithdrawMsgsResponse{
		Withdraws:  msgs,
		Pagination: pageRes,
	}, nil
}

// read data from kvstore for response of query request for params set
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}