package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Routes, Keys for liquidity module
const (
	// ModuleName is the name of the module.
	ModuleName = "liquidity"

	// RouterKey is the message route for the liquidity module.
	RouterKey = ModuleName

	// StoreKey is the default store key for the liquidity module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the liquidity module.
	QuerierRoute = ModuleName
)

// prefix key of liquidity states for indexing when kvstore
var (
	// param key for global Liquidity Pool IDs
	GlobalPoolIdKey = []byte("globalPoolId")

	PoolKeyPrefix               = []byte{0x11}
	PoolByReserveIndexKeyPrefix = []byte{0x12}

	PoolBatchIndexKeyPrefix = []byte{0x21} // LastPoolBatchIndex
	PoolBatchKeyPrefix      = []byte{0x22}

	PoolBatchDepositMsgIndexKeyPrefix  = []byte{0x31}
	PoolBatchWithdrawMsgIndexKeyPrefix = []byte{0x32}
	PoolBatchSwapMsgIndexKeyPrefix     = []byte{0x33}
)

// return kv indexing key of the pool
func GetPoolKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolId))
	return key
}

// return kv indexing key of the pool indexed by reserve account
func GetPoolByReserveAccIndexKey(reserveAcc sdk.AccAddress) []byte {
	return append(PoolByReserveIndexKeyPrefix, reserveAcc.Bytes()...)
}

// return kv indexing key of the latest index value of the pool batch
func GetPoolBatchIndexKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchIndexKeyPrefix[0]
	copy(key[1:], sdk.Uint64ToBigEndian(poolId))
	return key
}

// return kv indexing key of the pool batch indexed by pool id
func GetPoolBatchKey(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// Get prefix of the deposit batch messages that given pool for iteration
func GetPoolBatchDepositMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// Get prefix of the withdraw batch messages that given pool for iteration
func GetPoolBatchWithdrawMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// Get prefix of the swap batch messages that given pool for iteration
func GetPoolBatchSwapMsgsPrefix(poolId uint64) []byte {
	key := make([]byte, 9)
	key[0] = PoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	return key
}

// return kv indexing key of the latest index value of the msg index
func GetPoolBatchDepositMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = PoolBatchDepositMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

// return kv indexing key of the latest index value of the msg index
func GetPoolBatchWithdrawMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = PoolBatchWithdrawMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}

// return kv indexing key of the latest index value of the msg index
func GetPoolBatchSwapMsgIndexKey(poolId, msgIndex uint64) []byte {
	key := make([]byte, 17)
	key[0] = PoolBatchSwapMsgIndexKeyPrefix[0]
	copy(key[1:9], sdk.Uint64ToBigEndian(poolId))
	copy(key[9:17], sdk.Uint64ToBigEndian(msgIndex))
	return key
}