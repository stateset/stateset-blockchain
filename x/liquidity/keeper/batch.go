package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

// Reinitialize batch messages that were not executed in the previous batch and delete batch messages that were executed or ready to delete.
func (k Keeper) DeleteAndInitPoolBatch(ctx sdk.Context) {
	k.IterateAllPoolBatches(ctx, func(poolBatch types.PoolBatch) bool {
		// Delete and init next batch when not empty batch on before block
		if poolBatch.Executed {

			// On the other hand, BatchDeposit, BatchWithdraw, is all handled by the endblock if there is no error.
			// If there are BatchMsgs left, reset the Executed, Succeeded flag so that it can be executed in the next batch.
			depositMsgs := k.GetAllRemainingPoolBatchDepositMsgs(ctx, poolBatch)
			if len(depositMsgs) > 0 {
				for _, msg := range depositMsgs {
					msg.Executed = false
					msg.Succeeded = false
				}
				k.SetPoolBatchDepositMsgsByPointer(ctx, poolBatch.PoolId, depositMsgs)
			}

			withdrawMsgs := k.GetAllRemainingPoolBatchWithdrawMsgs(ctx, poolBatch)
			if len(withdrawMsgs) > 0 {
				for _, msg := range withdrawMsgs {
					msg.Executed = false
					msg.Succeeded = false
				}
				k.SetPoolBatchWithdrawMsgsByPointer(ctx, poolBatch.PoolId, withdrawMsgs)
			}

			height := ctx.BlockHeight()
			// reinitialize remaining batch msgs
			// In the case of BatchSwapMsgs, it is often fractional matched or has not yet expired since it has not passed ExpiryHeight.
			swapMsgs := k.GetAllRemainingPoolBatchSwapMsgs(ctx, poolBatch)
			if len(swapMsgs) > 0 {
				for _, msg := range swapMsgs {
					if height > msg.OrderExpiryHeight {
						msg.ToBeDeleted = true
					} else {
						msg.Executed = false
						msg.Succeeded = false
					}
				}
				k.SetPoolBatchSwapMsgPointers(ctx, poolBatch.PoolId, swapMsgs)
			}

			// delete batch messages that were executed or ready to delete
			k.DeleteAllReadyPoolBatchDepositMsgs(ctx, poolBatch)
			k.DeleteAllReadyPoolBatchWithdrawMsgs(ctx, poolBatch)
			k.DeleteAllReadyPoolBatchSwapMsgs(ctx, poolBatch)

			// Increase the batch index and initialize the values.
			k.InitNextBatch(ctx, poolBatch)
		}
		return false
	})
}

// Increase the index of the already executed batch for processing as the next batch and reinitialize the values.
func (k Keeper) InitNextBatch(ctx sdk.Context, poolBatch types.PoolBatch) error {
	if !poolBatch.Executed {
		return types.ErrBatchNotExecuted
	}
	poolBatch.BatchIndex = k.GetNextBatchIndexWithUpdate(ctx, poolBatch.PoolId)
	poolBatch.BeginHeight = ctx.BlockHeight()
	poolBatch.Executed = false
	k.SetPoolBatch(ctx, poolBatch)
	return nil
}

// In case of deposit, withdraw, and swap msgs, unlike other normal tx msgs,
// collect them in the liquidity pool batch and perform an execution once at the endblock to calculate and use the universal price.
func (k Keeper) ExecutePoolBatch(ctx sdk.Context) {
	k.IterateAllPoolBatches(ctx, func(poolBatch types.PoolBatch) bool {
		if !poolBatch.Executed {
			if poolBatch.Executed {
				return false
			}
			executedMsgCount, err := k.SwapExecution(ctx, poolBatch)
			if err != nil {
				panic(err)
			}
			k.IterateAllPoolBatchDepositMsgs(ctx, poolBatch, func(batchMsg types.BatchPoolDepositMsg) bool {
				executedMsgCount++
				if err := k.DepositPool(ctx, batchMsg); err != nil {
					k.RefundDepositPool(ctx, batchMsg)
				}
				return false
			})
			k.IterateAllPoolBatchWithdrawMsgs(ctx, poolBatch, func(batchMsg types.BatchPoolWithdrawMsg) bool {
				executedMsgCount++
				if err := k.WithdrawPool(ctx, batchMsg); err != nil {
					k.RefundWithdrawPool(ctx, batchMsg)
				}
				return false
			})
			// set executed when something executed
			if executedMsgCount > 0 {
				poolBatch.Executed = true
				k.SetPoolBatch(ctx, poolBatch)
			}
		}
		return false
	})
}

// In order to deal with the batch at once, the coins of msgs deposited in escrow.
func (k Keeper) HoldEscrow(ctx sdk.Context, depositor sdk.AccAddress, depositCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, depositCoins); err != nil {
		return err
	}
	return nil
}

// If batch messages has expired or has not been processed, will be refunded the escrow that had been deposited through this function.
func (k Keeper) ReleaseEscrow(ctx sdk.Context, withdrawer sdk.AccAddress, withdrawCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawer, withdrawCoins); err != nil {
		return err
	}
	return nil
}

// Generate inputs and outputs to treat escrow refunds atomically.
func (k Keeper) ReleaseEscrowForMultiSend(withdrawer sdk.AccAddress, withdrawCoins sdk.Coins) (
	banktypes.Input, banktypes.Output, error) {
	var input banktypes.Input
	var output banktypes.Output

	input = banktypes.NewInput(k.accountKeeper.GetModuleAddress(types.ModuleName), withdrawCoins)
	output = banktypes.NewOutput(withdrawer, withdrawCoins)

	if err := banktypes.ValidateInputsOutputs([]banktypes.Input{input}, []banktypes.Output{output}); err != nil {
		return banktypes.Input{}, banktypes.Output{}, err
	}
	return input, output, nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) DepositPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToPool) error {
	if err := k.ValidateMsgDepositPool(ctx, *msg); err != nil {
		return err
	}
	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolDepositMsg{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.DepositMsgIndex,
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetDepositor(), msg.DepositCoins); err != nil {
		return err
	}

	poolBatch.DepositMsgIndex += 1
	k.SetPoolBatch(ctx, poolBatch)
	k.SetPoolBatchDepositMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) WithdrawPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawFromPool) error {
	if err := k.ValidateMsgWithdrawPool(ctx, *msg); err != nil {
		return err
	}
	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return types.ErrPoolBatchNotExists
	}
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolWithdrawMsg{
		MsgHeight: ctx.BlockHeight(),
		MsgIndex:  poolBatch.WithdrawMsgIndex,
		Msg:       msg,
	}

	if err := k.HoldEscrow(ctx, msg.GetWithdrawer(), sdk.NewCoins(msg.PoolCoin)); err != nil {
		return err
	}

	poolBatch.WithdrawMsgIndex += 1
	k.SetPoolBatch(ctx, poolBatch)
	k.SetPoolBatchWithdrawMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return nil
}

// In order to deal with the batch at once, Put the message in the batch and the coins of the msgs deposited in escrow.
func (k Keeper) SwapPoolToBatch(ctx sdk.Context, msg *types.MsgSwap, OrderExpirySpanHeight int64) (*types.BatchPoolSwapMsg, error) {
	if err := k.ValidateMsgSwap(ctx, *msg); err != nil {
		return nil, err
	}
	poolBatch, found := k.GetPoolBatch(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolBatchNotExists
	}
	if poolBatch.BeginHeight == 0 {
		poolBatch.BeginHeight = ctx.BlockHeight()
	}

	batchPoolMsg := types.BatchPoolSwapMsg{
		MsgHeight:          ctx.BlockHeight(),
		MsgIndex:           poolBatch.SwapMsgIndex,
		Executed:           false,
		Succeeded:          false,
		ToBeDeleted:        false,
		ExchangedOfferCoin: sdk.NewCoin(msg.OfferCoin.Denom, sdk.ZeroInt()),
		RemainingOfferCoin: msg.OfferCoin,
		Msg:                msg,
	}
	// TODO: add logic if OrderExpiryHeight==0, pass on batch logic
	batchPoolMsg.OrderExpiryHeight = batchPoolMsg.MsgHeight + OrderExpirySpanHeight

	if err := k.HoldEscrow(ctx, msg.GetSwapRequester(), sdk.NewCoins(msg.OfferCoin)); err != nil {
		return nil, err
	}

	poolBatch.SwapMsgIndex += 1
	k.SetPoolBatch(ctx, poolBatch)
	k.SetPoolBatchSwapMsg(ctx, poolBatch.PoolId, batchPoolMsg)
	// TODO: msg event with msgServer after rebase stargate version sdk
	return &batchPoolMsg, nil
}
