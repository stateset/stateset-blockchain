package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the distribution MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Message server, handler for CreatePool msg
func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.Keeper.CreatePool(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeCreatePool,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.PoolCreatorAddress),
			sdk.NewAttribute(types.AttributeValuePoolId, ""),
			sdk.NewAttribute(types.AttributeValuePoolTypeIndex, fmt.Sprintf("%d", msg.PoolTypeIndex)),
			sdk.NewAttribute(types.AttributeValueReserveCoinDenoms, ""),
			sdk.NewAttribute(types.AttributeValueReserveAccount, ""),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, ""),
			sdk.NewAttribute(types.AttributeValueSwapFeeRate, ""),
			sdk.NewAttribute(types.AttributeValuePoolFeeRate, ""),
			sdk.NewAttribute(types.AttributeValueBatchSize, ""),
		),
	)
	return &types.MsgCreatePoolResponse{}, nil
}

// Message server, handler for MsgDepositToPool
func (k msgServer) DepositToPool(goCtx context.Context, msg *types.MsgDepositToPool) (*types.MsgDepositToPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.Keeper.DepositPoolToBatch(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeDepositToPoolToBatch,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DepositorAddress),
			sdk.NewAttribute(types.AttributeValueBatchID, ""),
		),
	)
	return &types.MsgDepositToPoolResponse{}, nil
}

// Message server, handler for MsgWithdrawFromPool
func (k msgServer) WithdrawFromPool(goCtx context.Context, msg *types.MsgWithdrawFromPool) (*types.MsgWithdrawFromPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.Keeper.WithdrawPoolToBatch(ctx, msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeWithdrrawFromPoolToBatch,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.WithdrawerAddress),
			sdk.NewAttribute(types.AttributeValueBatchID, ""),
		),
	)
	return &types.MsgWithdrawFromPoolResponse{}, nil
}

// Message server, handler for MsgSwap
func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if _, err := k.Keeper.SwapPoolToBatch(ctx, msg, 0); err != nil {
		return &types.MsgSwapResponse{}, err
	}
	return &types.MsgSwapResponse{}, nil
}