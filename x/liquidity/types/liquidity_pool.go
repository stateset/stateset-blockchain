package types

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (lp Pool) GetPoolKey() string {
	return GetPoolKey(lp.ReserveCoinDenoms, lp.PoolTypeIndex)
}

// Validate each constraint of the liquidity pool
func (lp Pool) Validate() error {
	if lp.PoolId == 0 {
		return ErrPoolNotExists
	}
	if lp.PoolTypeIndex == 0 {
		return ErrPoolTypeNotExists
	}
	if lp.ReserveCoinDenoms == nil || len(lp.ReserveCoinDenoms) == 0 {
		return ErrNumOfReserveCoinDenoms
	}
	if uint32(len(lp.ReserveCoinDenoms)) > MaxReserveCoinNum || uint32(len(lp.ReserveCoinDenoms)) < MinReserveCoinNum {
		return ErrNumOfReserveCoinDenoms
	}
	sortedDenomA, sortedDenomB := AlphabeticalDenomPair(lp.ReserveCoinDenoms[0], lp.ReserveCoinDenoms[1])
	if sortedDenomA != lp.ReserveCoinDenoms[0] || sortedDenomB != lp.ReserveCoinDenoms[1] {
		return ErrBadOrderingReserveCoinDenoms
	}
	if lp.ReserveAccountAddress == "" {
		return ErrEmptyReserveAccountAddress
	}
	//addr, err := sdk.AccAddressFromBech32(lp.ReserveAccountAddress)
	//if err != nil || lp.GetReserveAccount().Equals(addr) {
	//	return ErrBadReserveAccountAddress
	//}
	if lp.ReserveAccountAddress != GetPoolReserveAcc(lp.GetPoolKey()).String() {
		return ErrBadReserveAccountAddress
	}
	if lp.PoolCoinDenom == "" {
		return ErrEmptyPoolCoinDenom
	}
	if lp.PoolCoinDenom != lp.GetPoolKey() {
		return ErrBadPoolCoinDenom
	}
	return nil
}

// Calculate unique Pool key of the liquidity pool
func GetPoolKey(reserveCoinDenoms []string, poolTypeIndex uint32) string {
	return strings.Join(append(reserveCoinDenoms, strconv.FormatUint(uint64(poolTypeIndex), 10)), "-")
}

func NewPoolBatch(poolId, batchIndex uint64) PoolBatch {
	return PoolBatch{
		PoolId:           poolId,
		BatchIndex:       batchIndex,
		BeginHeight:      0,
		DepositMsgIndex:  1,
		WithdrawMsgIndex: 1,
		SwapMsgIndex:     1,
		Executed:         false,
	}
}

// MustMarshalPool returns the pool bytes. Panics if fails
func MustMarshalPool(cdc codec.BinaryMarshaler, pool Pool) []byte {
	return cdc.MustMarshalBinaryBare(&pool)
}

// MustUnmarshalPool return the unmarshaled pool from bytes.
// Panics if fails.
func MustUnmarshalPool(cdc codec.BinaryMarshaler, value []byte) Pool {
	pool, err := UnmarshalPool(cdc, value)
	if err != nil {
		panic(err)
	}

	return pool
}

// return the pool
func UnmarshalPool(cdc codec.BinaryMarshaler, value []byte) (pool Pool, err error) {
	err = cdc.UnmarshalBinaryBare(value, &pool)
	return pool, err
}

// return sdk.AccAddress object of he address saved as string because of protobuf
func (lp Pool) GetReserveAccount() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(lp.ReserveAccountAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// return pool coin denom of the liquidity poool
func (lp Pool) GetPoolCoinDenom() string { return lp.PoolCoinDenom }

// return pool id of the liquidity poool
func (lp Pool) GetPoolId() uint64 { return lp.PoolId }

// Pools is a collection of pools
type Pools []Pool

// PoolsBatch is a collection of poolBatch
type PoolsBatch []PoolBatch

// get string of list of liquidity pool
func (lps Pools) String() (out string) {
	for _, del := range lps {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// MustMarshalPoolBatch returns the PoolBatch bytes. Panics if fails
func MustMarshalPoolBatch(cdc codec.BinaryMarshaler, poolBatch PoolBatch) []byte {
	return cdc.MustMarshalBinaryBare(&poolBatch)
}

// return the poolBatch
func UnmarshalPoolBatch(cdc codec.BinaryMarshaler, value []byte) (poolBatch PoolBatch, err error) {
	err = cdc.UnmarshalBinaryBare(value, &poolBatch)
	return poolBatch, err
}

// MustUnmarshalPool return the unmarshaled PoolBatch from bytes.
// Panics if fails.
func MustUnmarshalPoolBatch(cdc codec.BinaryMarshaler, value []byte) PoolBatch {
	poolBatch, err := UnmarshalPoolBatch(cdc, value)
	if err != nil {
		panic(err)
	}

	return poolBatch
}

// MustMarshalBatchPoolDepositMsg returns the BatchPoolDepositMsg bytes. Panics if fails
func MustMarshalBatchPoolDepositMsg(cdc codec.BinaryMarshaler, msg BatchPoolDepositMsg) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the BatchPoolDepositMsg
func UnmarshalBatchPoolDepositMsg(cdc codec.BinaryMarshaler, value []byte) (msg BatchPoolDepositMsg, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalBatchPoolDepositMsg return the unmarshaled BatchPoolDepositMsg from bytes.
// Panics if fails.
func MustUnmarshalBatchPoolDepositMsg(cdc codec.BinaryMarshaler, value []byte) BatchPoolDepositMsg {
	msg, err := UnmarshalBatchPoolDepositMsg(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalBatchPoolWithdrawMsg returns the BatchPoolWithdrawMsg bytes. Panics if fails
func MustMarshalBatchPoolWithdrawMsg(cdc codec.BinaryMarshaler, msg BatchPoolWithdrawMsg) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the BatchPoolWithdrawMsg
func UnmarshalBatchPoolWithdrawMsg(cdc codec.BinaryMarshaler, value []byte) (msg BatchPoolWithdrawMsg, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalBatchPoolWithdrawMsg return the unmarshaled BatchPoolWithdrawMsg from bytes.
// Panics if fails.
func MustUnmarshalBatchPoolWithdrawMsg(cdc codec.BinaryMarshaler, value []byte) BatchPoolWithdrawMsg {
	msg, err := UnmarshalBatchPoolWithdrawMsg(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}

// MustMarshalBatchPoolSwapMsg returns the BatchPoolSwapMsg bytes. Panics if fails
func MustMarshalBatchPoolSwapMsg(cdc codec.BinaryMarshaler, msg BatchPoolSwapMsg) []byte {
	return cdc.MustMarshalBinaryBare(&msg)
}

// return the UnmarshalBatchPoolSwapMsg
func UnmarshalBatchPoolSwapMsg(cdc codec.BinaryMarshaler, value []byte) (msg BatchPoolSwapMsg, err error) {
	err = cdc.UnmarshalBinaryBare(value, &msg)
	return msg, err
}

// MustUnmarshalBatchPoolSwapMsg return the unmarshaled BatchPoolSwapMsg from bytes.
// Panics if fails.
func MustUnmarshalBatchPoolSwapMsg(cdc codec.BinaryMarshaler, value []byte) BatchPoolSwapMsg {
	msg, err := UnmarshalBatchPoolSwapMsg(cdc, value)
	if err != nil {
		panic(err)
	}
	return msg
}
