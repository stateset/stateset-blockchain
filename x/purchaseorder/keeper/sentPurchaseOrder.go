package keeper

import (
	"encoding/binary"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stateset/stateset-blockchain/x/purchaseorder/types"
)

// GetSentPurchaseOrderCount get the total number of sentPurchaseOrder
func (k Keeper) GetSentPurchaseOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderCountKey))
	byteKey := types.KeyPrefix(types.SentPurchaseOrderCountKey)
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

// SetSentPurchaseOrderCount set the total number of sentPurchaseOrder
func (k Keeper) SetSentPurchaseOrderCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderCountKey))
	byteKey := types.KeyPrefix(types.SentPurchaseOrderCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendSentPurchaseOrder appends a sentPurchaseOrder in the store with a new id and update the count
func (k Keeper) AppendSentPurchaseOrder(
	ctx sdk.Context,
	sentPurchaseOrder types.SentPurchaseOrder,
) uint64 {
	// Create the sentPurchaseOrder
	count := k.GetSentPurchaseOrderCount(ctx)

	// Set the ID of the appended value
	sentPurchaseOrder.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderKey))
	appendedValue := k.cdc.MustMarshalBinaryBare(&sentPurchaseOrder)
	store.Set(GetSentPurchaseOrderIDBytes(sentPurchaseOrder.Id), appendedValue)

	// Update sentPurchaseOrder count
	k.SetSentPurchaseOrderCount(ctx, count+1)

	return count
}

// SetSentPurchaseOrder set a specific sentPurchaseOrder in the store
func (k Keeper) SetSentPurchaseOrder(ctx sdk.Context, sentPurchaseOrder types.SentPurchaseOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderKey))
	b := k.cdc.MustMarshalBinaryBare(&sentPurchaseOrder)
	store.Set(GetSentPurchaseOrderIDBytes(sentPurchaseOrder.Id), b)
}

// GetSentPurchaseOrder returns a sentPurchaseOrder from its id
func (k Keeper) GetSentPurchaseOrder(ctx sdk.Context, id uint64) types.SentPurchaseOrder {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderKey))
	var sentPurchaseOrder types.SentPurchaseOrder
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetSentPurchaseOrderIDBytes(id)), &sentPurchaseOrder)
	return sentPurchaseOrder
}

// HasSentPurchaseOrder checks if the sentPurchaseOrder exists in the store
func (k Keeper) HasSentPurchaseOrder(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderKey))
	return store.Has(GetSentPurchaseOrderIDBytes(id))
}

// GetSentPurchaseOrderOwner returns the creator of the sentPurchaseOrder
func (k Keeper) GetSentPurchaseOrderOwner(ctx sdk.Context, id uint64) string {
	return k.GetSentPurchaseOrder(ctx, id).Creator
}

// RemoveSentPurchaseOrder removes a sentPurchaseOrder from the store
func (k Keeper) RemoveSentPurchaseOrder(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderKey))
	store.Delete(GetSentPurchaseOrderIDBytes(id))
}

// GetAllSentPurchaseOrder returns all sentPurchaseOrder
func (k Keeper) GetAllSentPurchaseOrder(ctx sdk.Context) (list []types.SentPurchaseOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SentPurchaseOrderKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SentPurchaseOrder
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetSentPurchaseOrderIDBytes returns the byte representation of the ID
func GetSentPurchaseOrderIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetSentPurchaseOrderIDFromBytes returns ID in uint64 format from a byte array
func GetSentPurchaseOrderIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
