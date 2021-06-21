package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/stateset/stateset-blockchain/x/liquidity/types"
)

// read form kvstore and return a specific pool
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (pool types.Pool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolKey(poolId)

	value := store.Get(key)
	if value == nil {
		return pool, false
	}

	pool = types.MustUnmarshalPool(k.cdc, value)

	return pool, true
}

// store to kvstore a specific pool
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalPool(k.cdc, pool)
	store.Set(types.GetPoolKey(pool.PoolId), b)
}

// delete from kvstore a specific pool
func (k Keeper) DeletePool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	Key := types.GetPoolKey(pool.PoolId)
	store.Delete(Key)
}

// IterateAllPools iterate through all of the pools
func (k Keeper) IterateAllPools(ctx sdk.Context, cb func(pool types.Pool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PoolKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pool := types.MustUnmarshalPool(k.cdc, iterator.Value())
		if cb(pool) {
			break
		}
	}
}

// GetAllPools returns all pools used during genesis dump
func (k Keeper) GetAllPools(ctx sdk.Context) (pools []types.Pool) {
	k.IterateAllPools(ctx, func(pool types.Pool) bool {
		pools = append(pools, pool)
		return false
	})

	return pools
}

// GetNextLiquidityID returns and increments the global Liquidity Pool ID counter.
// If the global account number is not set, it initializes it with value 0.
func (k Keeper) GetNextPoolIdWithUpdate(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	poolId := k.GetNextPoolId(ctx)
	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: poolId + 1})
	store.Set(types.GlobalPoolIdKey, bz)
	return poolId
}

// return next liquidity pool id for new pool, using index of latest pool id
func (k Keeper) GetNextPoolId(ctx sdk.Context) uint64 {
	var poolId uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GlobalPoolIdKey)
	if bz == nil {
		// initialize the PoolId
		poolId = 1
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.UnmarshalBinaryBare(bz, &val)
		if err != nil {
			panic(err)
		}

		poolId = val.GetValue()
	}
	return poolId
}

// read form kvstore and return a specific pool indexed by given reserve account
func (k Keeper) GetPoolByReserveAccIndex(ctx sdk.Context, reserveAcc sdk.AccAddress) (pool types.Pool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolByReserveAccIndexKey(reserveAcc)

	value := store.Get(key)
	if value == nil {
		return pool, false
	}

	pool = types.MustUnmarshalPool(k.cdc, value)

	return pool, true
}

// Set Index by ReserveAcc for liquidity Pool duplication check
func (k Keeper) SetPoolByReserveAccIndex(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalPool(k.cdc, pool)
	store.Set(types.GetPoolByReserveAccIndexKey(pool.GetReserveAccount()), b)
}

// Set Liquidity Pool with set global pool id index +1 and index by reserveAcc
func (k Keeper) SetPoolAtomic(ctx sdk.Context, pool types.Pool) types.Pool {
	pool.PoolId = k.GetNextPoolIdWithUpdate(ctx)
	k.SetPool(ctx, pool)
	k.SetPoolByReserveAccIndex(ctx, pool)
	return pool
}

// return a specific GetPoolBatchIndexKey
func (k Keeper) GetPoolBatchIndex(ctx sdk.Context, poolId uint64) (poolBatchIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolBatchIndexKey(poolId)

	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	poolBatchIndex = sdk.BigEndianToUint64(bz)
	return poolBatchIndex
}

// set index for liquidity pool batch, it should be increase after batch executed
func (k Keeper) SetPoolBatchIndex(ctx sdk.Context, poolId, batchIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	b := sdk.Uint64ToBigEndian(batchIndex)
	store.Set(types.GetPoolBatchIndexKey(poolId), b)
}

// return a specific poolBatch
func (k Keeper) GetPoolBatch(ctx sdk.Context, poolId uint64) (poolBatch types.PoolBatch, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolBatchKey(poolId)

	value := store.Get(key)
	if value == nil {
		return poolBatch, false
	}

	poolBatch = types.MustUnmarshalPoolBatch(k.cdc, value)

	return poolBatch, true
}

// return next batch index, with set index increased
func (k Keeper) GetNextBatchIndexWithUpdate(ctx sdk.Context, poolId uint64) (batchIndex uint64) {
	batchIndex = k.GetPoolBatchIndex(ctx, poolId)
	batchIndex += 1
	k.SetPoolBatchIndex(ctx, poolId, batchIndex)
	return
}

// Get All batches of the all existed liquidity pools
func (k Keeper) GetAllPoolBatches(ctx sdk.Context) (poolBatches []types.PoolBatch) {
	k.IterateAllPoolBatches(ctx, func(poolBatch types.PoolBatch) bool {
		poolBatches = append(poolBatches, poolBatch)
		return false
	})

	return poolBatches
}

// IterateAllPoolBatches iterate through all of the poolBatches
func (k Keeper) IterateAllPoolBatches(ctx sdk.Context, cb func(poolBatch types.PoolBatch) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PoolBatchKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolBatch := types.MustUnmarshalPoolBatch(k.cdc, iterator.Value())
		if cb(poolBatch) {
			break
		}
	}
}

// Delete batch of the liquidity pool, it used for test case
func (k Keeper) DeletePoolBatch(ctx sdk.Context, poolBatch types.PoolBatch) {
	store := ctx.KVStore(k.storeKey)
	batchKey := types.GetPoolBatchKey(poolBatch.PoolId)
	store.Delete(batchKey)
}

// set batch of the liquidity pool, with current state
func (k Keeper) SetPoolBatch(ctx sdk.Context, poolBatch types.PoolBatch) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalPoolBatch(k.cdc, poolBatch)
	store.Set(types.GetPoolBatchKey(poolBatch.PoolId), b)
}

// return a specific poolBatchDepositMsg
func (k Keeper) GetPoolBatchDepositMsg(ctx sdk.Context, poolId, msgIndex uint64) (msg types.BatchPoolDepositMsg, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolBatchDepositMsgIndexKey(poolId, msgIndex)

	value := store.Get(key)
	if value == nil {
		return msg, false
	}

	msg = types.MustUnmarshalBatchPoolDepositMsg(k.cdc, value)
	return msg, true
}

// set deposit batch msg of the liquidity pool batch, with current state
func (k Keeper) SetPoolBatchDepositMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolDepositMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolDepositMsg(k.cdc, msg)
	store.Set(types.GetPoolBatchDepositMsgIndexKey(poolId, msg.MsgIndex), b)
}

// set deposit batch msgs of the liquidity pool batch, with current state using pointers
func (k Keeper) SetPoolBatchDepositMsgsByPointer(ctx sdk.Context, poolId uint64, msgList []*types.BatchPoolDepositMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolDepositMsg(k.cdc, *msg)
		store.Set(types.GetPoolBatchDepositMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

// set deposit batch msgs of the liquidity pool batch, with current state
func (k Keeper) SetPoolBatchDepositMsgs(ctx sdk.Context, poolId uint64, msgList []types.BatchPoolDepositMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolDepositMsg(k.cdc, msg)
		store.Set(types.GetPoolBatchDepositMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

//func (k Keeper) DeletePoolBatchDepositMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
//	store := ctx.KVStore(k.storeKey)
//	batchKey := types.GetPoolBatchDepositMsgIndexKey(poolId, msgIndex)
//	store.Delete(batchKey)
//}

// IterateAllPoolBatchDepositMsgs iterate through all of the PoolBatchDepositMsgs
func (k Keeper) IterateAllPoolBatchDepositMsgs(ctx sdk.Context, poolBatch types.PoolBatch, cb func(msg types.BatchPoolDepositMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetPoolBatchDepositMsgsPrefix(poolBatch.PoolId)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolDepositMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// IterateAllBatchDepositMsgs iterate through all of the BatchDepositMsgs of all batches
func (k Keeper) IterateAllBatchDepositMsgs(ctx sdk.Context, cb func(msg types.BatchPoolDepositMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.PoolBatchDepositMsgIndexKeyPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolDepositMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllBatchDepositMsgs returns all BatchDepositMsgs for all batches.
func (k Keeper) GetAllBatchDepositMsgs(ctx sdk.Context) (msgs []types.BatchPoolDepositMsg) {
	k.IterateAllBatchDepositMsgs(ctx, func(msg types.BatchPoolDepositMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllPoolBatchDepositMsgs returns all BatchDepositMsgs indexed by the poolBatch
func (k Keeper) GetAllPoolBatchDepositMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []types.BatchPoolDepositMsg) {
	k.IterateAllPoolBatchDepositMsgs(ctx, poolBatch, func(msg types.BatchPoolDepositMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllToDeletePoolBatchDepositMsgs returns all Not toDelete BatchDepositMsgs indexed by the poolBatch
func (k Keeper) GetAllNotToDeletePoolBatchDepositMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []types.BatchPoolDepositMsg) {
	k.IterateAllPoolBatchDepositMsgs(ctx, poolBatch, func(msg types.BatchPoolDepositMsg) bool {
		if !msg.ToBeDeleted {
			msgs = append(msgs, msg)
		}
		return false
	})
	return msgs
}

// GetAllRemainingPoolBatchDepositMsgs returns All only remaining BatchDepositMsgs after endblock , executed but not toDelete
func (k Keeper) GetAllRemainingPoolBatchDepositMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []*types.BatchPoolDepositMsg) {
	k.IterateAllPoolBatchDepositMsgs(ctx, poolBatch, func(msg types.BatchPoolDepositMsg) bool {
		if msg.Executed && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// delete deposit batch msgs of the liquidity pool batch which has state ToBeDeleted
func (k Keeper) DeleteAllReadyPoolBatchDepositMsgs(ctx sdk.Context, poolBatch types.PoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPoolBatchDepositMsgsPrefix(poolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolDepositMsg(k.cdc, iterator.Value())
		if msg.ToBeDeleted {
			store.Delete(iterator.Key())
		}
	}
}

// return a specific poolBatchWithdrawMsg
func (k Keeper) GetPoolBatchWithdrawMsg(ctx sdk.Context, poolId, msgIndex uint64) (msg types.BatchPoolWithdrawMsg, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolBatchWithdrawMsgIndexKey(poolId, msgIndex)

	value := store.Get(key)
	if value == nil {
		return msg, false
	}

	msg = types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, value)
	return msg, true
}

// set withdraw batch msg of the liquidity pool batch, with current state
func (k Keeper) SetPoolBatchWithdrawMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolWithdrawMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, msg)
	store.Set(types.GetPoolBatchWithdrawMsgIndexKey(poolId, msg.MsgIndex), b)
}

// set withdraw batch msgs of the liquidity pool batch, with current state using pointers
func (k Keeper) SetPoolBatchWithdrawMsgsByPointer(ctx sdk.Context, poolId uint64, msgList []*types.BatchPoolWithdrawMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, *msg)
		store.Set(types.GetPoolBatchWithdrawMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

// set withdraw batch msgs of the liquidity pool batch, with current state
func (k Keeper) SetPoolBatchWithdrawMsgs(ctx sdk.Context, poolId uint64, msgList []types.BatchPoolWithdrawMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolWithdrawMsg(k.cdc, msg)
		store.Set(types.GetPoolBatchWithdrawMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

//func (k Keeper) DeletePoolBatchWithdrawMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
//	store := ctx.KVStore(k.storeKey)
//	batchKey := types.GetPoolBatchWithdrawMsgIndexKey(poolId, msgIndex)
//	store.Delete(batchKey)
//}

// IterateAllPoolBatchWithdrawMsgs iterate through all of the PoolBatchWithdrawMsgs
func (k Keeper) IterateAllPoolBatchWithdrawMsgs(ctx sdk.Context, poolBatch types.PoolBatch, cb func(msg types.BatchPoolWithdrawMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetPoolBatchWithdrawMsgsPrefix(poolBatch.PoolId)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// IterateAllBatchWithdrawMsgs iterate through all of the BatchPoolWithdrawMsg of all batches
func (k Keeper) IterateAllBatchWithdrawMsgs(ctx sdk.Context, cb func(msg types.BatchPoolWithdrawMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.PoolBatchWithdrawMsgIndexKeyPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllBatchWithdrawMsgs returns all BatchWithdrawMsgs for all batches
func (k Keeper) GetAllBatchWithdrawMsgs(ctx sdk.Context) (msgs []types.BatchPoolWithdrawMsg) {
	k.IterateAllBatchWithdrawMsgs(ctx, func(msg types.BatchPoolWithdrawMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllPoolBatchWithdrawMsgs returns all BatchWithdrawMsgs indexed by the poolBatch
func (k Keeper) GetAllPoolBatchWithdrawMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []types.BatchPoolWithdrawMsg) {
	k.IterateAllPoolBatchWithdrawMsgs(ctx, poolBatch, func(msg types.BatchPoolWithdrawMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllToDeletePoolBatchWithdrawMsgs returns all Not to delete BatchWithdrawMsgs indexed by the poolBatch
func (k Keeper) GetAllNotToDeletePoolBatchWithdrawMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []types.BatchPoolWithdrawMsg) {
	k.IterateAllPoolBatchWithdrawMsgs(ctx, poolBatch, func(msg types.BatchPoolWithdrawMsg) bool {
		if !msg.ToBeDeleted {
			msgs = append(msgs, msg)
		}
		return false
	})
	return msgs
}

// GetAllRemainingPoolBatchWithdrawMsgs returns All only remaining BatchWithdrawMsgs after endblock, executed but not toDelete
func (k Keeper) GetAllRemainingPoolBatchWithdrawMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []*types.BatchPoolWithdrawMsg) {
	k.IterateAllPoolBatchWithdrawMsgs(ctx, poolBatch, func(msg types.BatchPoolWithdrawMsg) bool {
		if msg.Executed && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// delete withdraw batch msgs of the liquidity pool batch which has state ToBeDeleted
func (k Keeper) DeleteAllReadyPoolBatchWithdrawMsgs(ctx sdk.Context, poolBatch types.PoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPoolBatchWithdrawMsgsPrefix(poolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolWithdrawMsg(k.cdc, iterator.Value())
		if msg.ToBeDeleted {
			store.Delete(iterator.Key())
		}
	}
}

// return a specific GetPoolBatchSwapMsg, not used currently
//func (k Keeper) GetPoolBatchSwapMsg(ctx sdk.Context, poolId, msgIndex uint64) (msg types.BatchPoolSwapMsg, found bool) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetPoolBatchSwapMsgIndexKey(poolId, msgIndex)
//
//	value := store.Get(key)
//	if value == nil {
//		return msg, false
//	}
//
//	msg = types.MustUnmarshalBatchPoolSwapMsg(k.cdc, value)
//	return msg, true
//}

// set swap batch msg of the liquidity pool batch, with current state
func (k Keeper) SetPoolBatchSwapMsg(ctx sdk.Context, poolId uint64, msg types.BatchPoolSwapMsg) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalBatchPoolSwapMsg(k.cdc, msg)
	store.Set(types.GetPoolBatchSwapMsgIndexKey(poolId, msg.MsgIndex), b)
}

// Delete swap batch msg of the liquidity pool batch, it used for test case
func (k Keeper) DeletePoolBatchSwapMsg(ctx sdk.Context, poolId uint64, msgIndex uint64) {
	store := ctx.KVStore(k.storeKey)
	batchKey := types.GetPoolBatchSwapMsgIndexKey(poolId, msgIndex)
	store.Delete(batchKey)
}

// IterateAllPoolBatchSwapMsgs iterate through all of the PoolBatchSwapMsgs
func (k Keeper) IterateAllPoolBatchSwapMsgs(ctx sdk.Context, poolBatch types.PoolBatch, cb func(msg types.BatchPoolSwapMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.GetPoolBatchSwapMsgsPrefix(poolBatch.PoolId)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// IterateAllBatchSwapMsgs iterate through all of the BatchPoolSwapMsg of all batches
func (k Keeper) IterateAllBatchSwapMsgs(ctx sdk.Context, cb func(msg types.BatchPoolSwapMsg) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	prefix := types.PoolBatchSwapMsgIndexKeyPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if cb(msg) {
			break
		}
	}
}

// GetAllBatchSwapMsgs returns all BatchSwapMsgs of all batches
func (k Keeper) GetAllBatchSwapMsgs(ctx sdk.Context) (msgs []types.BatchPoolSwapMsg) {
	k.IterateAllBatchSwapMsgs(ctx, func(msg types.BatchPoolSwapMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// delete swap batch msgs of the liquidity pool batch which has state ToBeDeleted
func (k Keeper) DeleteAllReadyPoolBatchSwapMsgs(ctx sdk.Context, poolBatch types.PoolBatch) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetPoolBatchSwapMsgsPrefix(poolBatch.PoolId))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		msg := types.MustUnmarshalBatchPoolSwapMsg(k.cdc, iterator.Value())
		if msg.ToBeDeleted {
			store.Delete(iterator.Key())
		}
	}
}

// GetAllPoolBatchSwapMsgsAsPointer returns all BatchSwapMsgs pointer indexed by the poolBatch
func (k Keeper) GetAllPoolBatchSwapMsgsAsPointer(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllPoolBatchSwapMsgs(ctx, poolBatch, func(msg types.BatchPoolSwapMsg) bool {
		msgs = append(msgs, &msg)
		return false
	})
	return msgs
}

// GetAllPoolBatchSwapMsgs returns all BatchSwapMsgs indexed by the poolBatch
func (k Keeper) GetAllPoolBatchSwapMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []types.BatchPoolSwapMsg) {
	k.IterateAllPoolBatchSwapMsgs(ctx, poolBatch, func(msg types.BatchPoolSwapMsg) bool {
		msgs = append(msgs, msg)
		return false
	})
	return msgs
}

// GetAllNotProcessedPoolBatchSwapMsgs returns All only not processed swap msgs, not executed with not succeed and not toDelete BatchSwapMsgs indexed by the poolBatch
func (k Keeper) GetAllNotProcessedPoolBatchSwapMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllPoolBatchSwapMsgs(ctx, poolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if !msg.Executed && !msg.Succeeded && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// GetAllRemainingPoolBatchSwapMsgs returns All only remaining after endblock swap msgs, executed but not toDelete
func (k Keeper) GetAllRemainingPoolBatchSwapMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllPoolBatchSwapMsgs(ctx, poolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if msg.Executed && !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// GetAllNotToDeletePoolBatchSwapMsgs returns All only not to delete swap msgs
func (k Keeper) GetAllNotToDeletePoolBatchSwapMsgs(ctx sdk.Context, poolBatch types.PoolBatch) (msgs []*types.BatchPoolSwapMsg) {
	k.IterateAllPoolBatchSwapMsgs(ctx, poolBatch, func(msg types.BatchPoolSwapMsg) bool {
		if !msg.ToBeDeleted {
			msgs = append(msgs, &msg)
		}
		return false
	})
	return msgs
}

// set swap batch msgs of the liquidity pool batch, with current state using pointers
func (k Keeper) SetPoolBatchSwapMsgPointers(ctx sdk.Context, poolId uint64, msgList []*types.BatchPoolSwapMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolSwapMsg(k.cdc, *msg)
		store.Set(types.GetPoolBatchSwapMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}

// set swap batch msgs of the liquidity pool batch, with current state
func (k Keeper) SetPoolBatchSwapMsgs(ctx sdk.Context, poolId uint64, msgList []types.BatchPoolSwapMsg) {
	for _, msg := range msgList {
		if poolId != msg.Msg.PoolId {
			continue
		}
		store := ctx.KVStore(k.storeKey)
		b := types.MustMarshalBatchPoolSwapMsg(k.cdc, msg)
		store.Set(types.GetPoolBatchSwapMsgIndexKey(poolId, msg.MsgIndex), b)
	}
}