package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

// GetTimedoutPurchaseOrderCount get the total number of timedoutPurchaseOrder
func (k Keeper) GetTimedoutPurchaseOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderCountKey))
	byteKey := types.KeyPrefix(types.TimedoutPurchaseOrderCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetTimedoutPurchaseOrderCount set the total number of timedoutPurchaseOrder
func (k Keeper) SetTimedoutPurchaseOrderCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderCountKey))
	byteKey := types.KeyPrefix(types.TimedoutPurchaseOrderCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendTimedoutPurchaseOrder appends a timedoutPurchaseOrder in the store with a new id and update the count
func (k Keeper) AppendTimedoutPurchaseOrder(
	ctx sdk.Context,
	timedoutPurchaseOrder types.TimedoutPurchaseOrder,
) uint64 {
	// Create the timedoutPurchaseOrder
	count := k.GetTimedoutPurchaseOrderCount(ctx)

	// Set the ID of the appended value
	timedoutPurchaseOrder.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderKey))
	appendedValue := k.cdc.MustMarshalBinaryBare(&timedoutPurchaseOrder)
	store.Set(GetTimedoutPurchaseOrderIDBytes(timedoutPurchaseOrder.Id), appendedValue)

	// Update timedoutPurchaseOrder count
	k.SetTimedoutPurchaseOrderCount(ctx, count+1)

	return count
}

// SetTimedoutPurchaseOrder set a specific timedoutPurchaseOrder in the store
func (k Keeper) SetTimedoutPurchaseOrder(ctx sdk.Context, timedoutPurchaseOrder types.TimedoutPurchaseOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderKey))
	b := k.cdc.MustMarshalBinaryBare(&timedoutPurchaseOrder)
	store.Set(GetTimedoutPurchaseOrderIDBytes(timedoutPurchaseOrder.Id), b)
}

// GetTimedoutPurchaseOrder returns a timedoutPurchaseOrder from its id
func (k Keeper) GetTimedoutPurchaseOrder(ctx sdk.Context, id uint64) types.TimedoutPurchaseOrder {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderKey))
	var timedoutPurchaseOrder types.TimedoutPurchaseOrder
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetTimedoutPurchaseOrderIDBytes(id)), &timedoutPurchaseOrder)
	return timedoutPurchaseOrder
}

// HasTimedoutPurchaseOrder checks if the timedoutPurchaseOrder exists in the store
func (k Keeper) HasTimedoutPurchaseOrder(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderKey))
	return store.Has(GetTimedoutPurchaseOrderIDBytes(id))
}

// GetTimedoutPurchaseOrderOwner returns the creator of the timedoutPurchaseOrder
func (k Keeper) GetTimedoutPurchaseOrderOwner(ctx sdk.Context, id uint64) string {
	return k.GetTimedoutPurchaseOrder(ctx, id).Creator
}

// RemoveTimedoutPurchaseOrder removes a timedoutPurchaseOrder from the store
func (k Keeper) RemoveTimedoutPurchaseOrder(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderKey))
	store.Delete(GetTimedoutPurchaseOrderIDBytes(id))
}

// GetAllTimedoutPurchaseOrder returns all timedoutPurchaseOrder
func (k Keeper) GetAllTimedoutPurchaseOrder(ctx sdk.Context) (list []types.TimedoutPurchaseOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TimedoutPurchaseOrderKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TimedoutPurchaseOrder
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetTimedoutPurchaseOrderIDBytes returns the byte representation of the ID
func GetTimedoutPurchaseOrderIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetTimedoutPurchaseOrderIDFromBytes returns ID in uint64 format from a byte array
func GetTimedoutPurchaseOrderIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
